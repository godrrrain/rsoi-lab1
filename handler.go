package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PersonResponse struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Work    string `json:"work"`
	Age     int    `json:"age"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	storage Storage
}

func NewHandler(storage Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) GetPerson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("failed to convert id param to int: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	person, err := h.storage.Get(id)

	if err != nil {
		fmt.Printf("failed to get person %s\n", err.Error())

		if err.Error() == "person not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PersonToResponse(person))
}

func (h *Handler) GetPersons(c *gin.Context) {
	persons := h.storage.GetAll()

	c.JSON(http.StatusOK, PersonsToResponse(persons))
}

func (h *Handler) CreatePerson(c *gin.Context) {
	var person Person

	if err := c.BindJSON(&person); err != nil {
		fmt.Printf("failed to bind Person: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	h.storage.Insert(&person)

	c.Header("Location", "/api/v1/persons/"+strconv.Itoa(person.ID))
	c.Status(http.StatusCreated)
}

func (h *Handler) UpdatePerson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("failed to convert id param to int: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var person Person

	if err := c.BindJSON(&person); err != nil {
		fmt.Printf("failed to bind person: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	person.ID = id

	err = h.storage.Update(&person)

	if err != nil {
		if err.Error() == "person not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PersonToResponse(person))
}

func (h *Handler) DeletePerson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("failed to convert id param to int: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	err = h.storage.Delete(id)
	if err != nil {
		fmt.Printf("failed to delete person %s\n", err.Error())

		if err.Error() == "person not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, MessageResponse{
		Message: "person deleted",
	})
}

func PersonToResponse(person Person) PersonResponse {
	return PersonResponse{
		Id:      person.ID,
		Name:    person.Name,
		Address: person.Address,
		Work:    person.Work,
		Age:     person.Age,
	}
}

func PersonsToResponse(persons []Person) []PersonResponse {
	if persons == nil {
		return nil
	}

	res := make([]PersonResponse, len(persons))

	for index, value := range persons {
		res[index] = PersonToResponse(value)
	}

	return res
}
