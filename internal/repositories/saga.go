package repositories

import (
	"fmt"
	"go-api-find-my-friend/internal/models"
	"go-api-find-my-friend/pkg/errors"
	"go-api-find-my-friend/pkg/storage_provider"
	"mime/multipart"

	"gorm.io/gorm"
)

type SagaStep interface {
	Execute() error
	Compensate() error
	GetName() string
	SetNext(next SagaStep)
	SetPrevious(prev SagaStep)
	GetNext() SagaStep
	GetPrevious() SagaStep
	IsExecuted() bool
	SetExecuted(executed bool)
}

type SagaOrchestrator struct {
	head SagaStep
	tail SagaStep
}

func NewSagaOrchestrator() *SagaOrchestrator {
	return &SagaOrchestrator{
		head: nil,
		tail: nil,
	}
}

func (s *SagaOrchestrator) AddSteps(steps ...SagaStep) {
	for _, step := range steps {
		s.addStep(step)
	}
}

func (s *SagaOrchestrator) addStep(step SagaStep) {
	if s.head == nil {
		s.head = step
		s.tail = step
		return
	}

	s.tail.SetNext(step)
	step.SetPrevious(s.tail)
	s.tail = step
}

func (s *SagaOrchestrator) Run() error {
	current := s.head

	for current != nil {
		if err := current.Execute(); err != nil {
			s.rollback(current)
			return err
		}
		current.SetExecuted(true)
		current = current.GetNext()
	}

	return nil
}

func (s *SagaOrchestrator) rollback(failedStep SagaStep) {
	current := failedStep.GetPrevious()

	for current != nil {
		if current.IsExecuted() {
			if err := current.Compensate(); err != nil {
				fmt.Printf("Warning: Failed to compensate step %s: %v\n", current.GetName(), err)
			}
		}
		current = current.GetPrevious()
	}
}

type UploadPictureStep struct {
	pet             *models.Pet
	picture         *multipart.FileHeader
	storageProvider storage_provider.StorageProvider
	pictureURL      string
	uploaded        bool
	executed        bool
	next            SagaStep
	previous        SagaStep
}

func NewUploadPictureStep(pet *models.Pet, picture *multipart.FileHeader, storageProvider storage_provider.StorageProvider) *UploadPictureStep {
	return &UploadPictureStep{
		pet:             pet,
		picture:         picture,
		storageProvider: storageProvider,
	}
}

func (s *UploadPictureStep) Execute() error {
	pictureURL, err := s.storageProvider.Upload(s.picture)
	if err != nil {
		return errors.NewInternalServerError("Failed to upload picture")
	}

	s.uploaded = true
	s.pet.PictureURL = pictureURL
	return nil
}

func (s *UploadPictureStep) Compensate() error {
	if s.uploaded && s.pictureURL != "" {
		return s.storageProvider.Delete(s.pictureURL)
	}
	return nil
}

func (s *UploadPictureStep) GetName() string {
	return "UploadPicture"
}

func (s *UploadPictureStep) GetPictureURL() string {
	return s.pictureURL
}

func (s *UploadPictureStep) SetNext(next SagaStep) {
	s.next = next
}

func (s *UploadPictureStep) SetPrevious(prev SagaStep) {
	s.previous = prev
}

func (s *UploadPictureStep) GetNext() SagaStep {
	return s.next
}

func (s *UploadPictureStep) GetPrevious() SagaStep {
	return s.previous
}

func (s *UploadPictureStep) IsExecuted() bool {
	return s.executed
}

func (s *UploadPictureStep) SetExecuted(executed bool) {
	s.executed = executed
}

type CreatePetStep struct {
	pet      *models.Pet
	db       *gorm.DB
	created  bool
	executed bool
	next     SagaStep
	previous SagaStep
}

func NewCreatePetStep(pet *models.Pet, db *gorm.DB) *CreatePetStep {
	return &CreatePetStep{
		pet: pet,
		db:  db,
	}
}

func (s *CreatePetStep) Execute() error {
	err := s.db.Create(s.pet).Error
	if err != nil {
		return errors.NewInternalServerError("Failed to create pet")
	}

	s.created = true
	return nil
}

func (s *CreatePetStep) Compensate() error {
	return nil
}

func (s *CreatePetStep) GetName() string {
	return "CreatePet"
}

func (s *CreatePetStep) SetNext(next SagaStep) {
	s.next = next
}

func (s *CreatePetStep) SetPrevious(prev SagaStep) {
	s.previous = prev
}

func (s *CreatePetStep) GetNext() SagaStep {
	return s.next
}

func (s *CreatePetStep) GetPrevious() SagaStep {
	return s.previous
}

func (s *CreatePetStep) IsExecuted() bool {
	return s.executed
}

func (s *CreatePetStep) SetExecuted(executed bool) {
	s.executed = executed
}
