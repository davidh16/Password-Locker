package service

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"password-lock/models"
)

func (s Service) EncryptPassword(secretKey string, password string) string {

	c, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())

	encryptedPassowrd := gcm.Seal(nonce, nonce, []byte(password), nil)

	return base64.StdEncoding.EncodeToString(encryptedPassowrd)
}

func (s Service) GetEntityIconPath(entityType int) string {
	switch entityType {
	case 0:
		return "/Users/davidhorvat/GolandProjects/Password-lock/logos/github.png"
	case 1:
		return "/Users/davidhorvat/GolandProjects/Password-lock/logos/facebook.png"
	case 2:
		return "/Users/davidhorvat/GolandProjects/Password-lock/logos/gmail.png"
	case 3:
		return "/Users/davidhorvat/GolandProjects/Password-lock/logos/linkedin.png"
	case 4:
		return "/Users/davidhorvat/GolandProjects/Password-lock/logos/instagram.png"
	}
	return ""
}

func (s Service) CreateEntity(entity models.Entity) (*models.Entity, error) {
	result := s.entityRepository.Db().Create(entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}

func (s Service) DeleteEntity(entityUuid string) error {
	result := s.entityRepository.Db().Where("uuid=?", entityUuid).Delete(models.Entity{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s Service) GetEntityByUuid(ctx *gin.Context, entityUuid string, secretKey string) (*models.Entity, error) {

	loggedInUser := s.Me(ctx)

	var entity models.Entity
	result := s.entityRepository.Db().Where("uuid=? AND user_uuid=?", entityUuid, loggedInUser).First(&entity)
	if result.Error != nil {
		return nil, result.Error
	}

	decryptedPassword, err := decryptEntityPassword(entity.Password, secretKey)
	if err != nil {
		return nil, err
	}

	entity.Password = decryptedPassword

	return &entity, nil
}

func decryptEntityPassword(password string, secretKey string) (string, error) {

	ciphertext, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		return "", err
	}

	key := []byte(secretKey)

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	decryptedPassword, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(decryptedPassword), nil

}
