package server

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// attributes of each user
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// array of our users
var users []User

// Login request and response JSON properties
type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	User        User   `json:"user"`
}

func (server *Server) login(ctx *gin.Context) {
	// Request binding for login credentials
	var req loginRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate credentials and generate an access token
	for _, user := range users {
		if user.Username == req.Username {
			if user.Password == req.Password {

				// Create and send an access token
				accessToken, err := server.tokenMaker.CreateToken(req.Username, time.Minute)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				// respond with login details
				resp := loginResponse{
					AccessToken: accessToken,
					User:        user,
				}

				ctx.JSON(http.StatusOK, resp)
				return
			}
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Incorrect password"})
			return
		}
	}

	ctx.JSON(http.StatusNotFound, users)
	return
}

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	// Request binding for new user details
	var user User
	err := ctx.ShouldBindJSON(&user) // method will populate the struct with the JSON data provided by the request body
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assign a unique ID and add the user to the list
	user.ID = strconv.Itoa(rand.Intn(1000))
	users = append(users, user)

	// Respond with te updated user list
	ctx.JSON(http.StatusOK, users)
	return
}

type deleteUserRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	// Request binding for user ID
	var req deleteUserRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find and delete the user fro mthe list
	for idx, user := range users {
		if user.ID == req.ID {
			// remove user from the users slice
			users = append(users[:idx], users[idx+1:]...)
			ctx.JSON(http.StatusOK, users)
			return
		}
	}

	// Respond with user list if the user is not found
	ctx.JSON(http.StatusNotFound, users)
	return
}
