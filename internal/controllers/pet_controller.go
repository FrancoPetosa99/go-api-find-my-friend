package controllers

import (
	"net/http"
	"strconv"
	"sync"

	"go-api-find-my-friend/internal/services"
	"go-api-find-my-friend/pkg/errors"
	"go-api-find-my-friend/pkg/pagination"

	"github.com/gin-gonic/gin"
)

var (
	petControllerInstance *PetController
	petControllerOnce     sync.Once
)

var (
	ErrInvalidPetID         = errors.NewBadRequestError("invalid pet ID")
	ErrCreatePetInvalidBody = errors.NewBadRequestError("invalid body")
	ErrInvalidQueryParams   = errors.NewBadRequestError("invalid query params")
	ErrUpdatePetInvalidBody = errors.NewBadRequestError("invalid body")
)

type PetController struct {
	petService  *services.PetService
	userService *services.UserService
}

func NewPetController() *PetController {
	petControllerOnce.Do(func() {
		petControllerInstance = &PetController{
			petService:  services.NewPetService(),
			userService: services.NewUserService(),
		}
	})
	return petControllerInstance
}

func (c *PetController) CreatePet(ctx *gin.Context) {
	var dto services.PetCreateDTO

	if err := ctx.ShouldBind(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrCreatePetInvalidBody)
		return
	}

	errs := map[string]string{}
	passed := dto.Validate(&errs)
	if !passed {
		ctx.JSON(http.StatusBadRequest, errs)
		return
	}

	userID, _ := ctx.Get("user_id")
	pet, err := c.petService.CreatePet(&dto, userID.(int))
	if err != nil {
		ctx.JSON(getErrStatusCode(err), err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Pet created successfully",
		"pet": gin.H{
			"id":          pet.ID,
			"name":        pet.Name,
			"description": pet.Description,
			"type":        pet.Type,
			"breed":       pet.Breed,
			"picture_url": pet.PictureURL,
		},
	})
}

func (c *PetController) SearchPets(ctx *gin.Context) {
	var dto SearchPetsPaginationDTO

	if err := ctx.ShouldBindQuery(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrInvalidQueryParams)
		return
	}

	searchParams := pagination.PaginationParams{
		Page:    dto.Page,
		Size:    dto.Size,
		SortDir: dto.SortDir,
	}

	filterParams := pagination.FilterPet{}

	if dto.Type != "" {
		filterParams.Type = &dto.Type
	}
	if dto.Breed != "" {
		filterParams.Breed = &dto.Breed
	}
	if dto.LastSeenPlace != "" {
		filterParams.LastSeenPlace = &dto.LastSeenPlace
	}

	pagination, err := c.petService.SearchPets(&filterParams, &searchParams)
	if err != nil {
		ctx.JSON(getErrStatusCode(err), err)
		return
	}

	ctx.JSON(http.StatusOK, pagination)
}

func (c *PetController) GetPet(ctx *gin.Context) {
	petID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrInvalidPetID)
		return
	}

	pet, err := c.petService.GetPetByID(petID)
	if err != nil {
		ctx.JSON(getErrStatusCode(err), err)
		return
	}

	userID, _ := ctx.Get("user_id")
	isOwner := pet.UserID == userID

	ctx.JSON(http.StatusOK,
		PetDetailDTO{
			PetID:         pet.ID,
			OwnerID:       pet.UserID,
			OwnerName:     pet.User.Name,
			OwnerLastName: pet.User.LastName,
			OwnerEmail:    pet.User.Email,
			OwnerPhone:    pet.User.Phone,
			Name:          pet.Name,
			Description:   pet.Description,
			Type:          pet.Type,
			Breed:         pet.Breed,
			LastSeenTime:  pet.LastSeenTime,
			LastSeenPlace: pet.LastSeenPlace,
			PictureURL:    pet.PictureURL,
			IsFound:       pet.IsFound,
			CanEdit:       isOwner,
			CanDelete:     isOwner,
		},
	)
}

func (c *PetController) UpdatePet(ctx *gin.Context) {
	petID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrInvalidPetID)
		return
	}

	var dto services.PetUpdateDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrUpdatePetInvalidBody)
		return
	}

	errs := map[string]string{}
	passed := dto.Validate(&errs)
	if !passed {
		ctx.JSON(http.StatusBadRequest, errs)
		return
	}

	userID, _ := ctx.Get("user_id")
	err = c.petService.UpdatePet(userID.(int), petID, &dto)
	if err != nil {
		ctx.JSON(getErrStatusCode(err), err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *PetController) UpdatePetMarkAsFound(ctx *gin.Context) {
	petID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrInvalidPetID)
		return
	}

	userID, _ := ctx.Get("user_id")
	err = c.petService.UpdatePetAsFound(userID.(int), petID)
	if err != nil {
		ctx.JSON(getErrStatusCode(err), err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *PetController) DeletePet(ctx *gin.Context) {
	petID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrInvalidPetID)
		return
	}

	userID, _ := ctx.Get("user_id")
	err = c.petService.DeletePet(userID.(int), petID)
	if err != nil {
		ctx.JSON(getErrStatusCode(err), err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
