package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
	"github.com/nazzarr03/social-media-restful-api/controller"
	"github.com/nazzarr03/social-media-restful-api/middleware"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := fiber.New()

	secret := os.Getenv("JWT_SECRET")
	r.Use(middleware.Authentication(secret))

	r.Post("/login", controller.Login)
	r.Post("/signup", controller.SignUp)
	r.Post("/logout", controller.Logout)
	r.Post("/profilepicture/:id", controller.UploadProfilePicture)
	r.Post("/updateprofilepicture/:id", controller.UpdateProfilePicture)

	r.Post("/createpost/:id", controller.CreatePost)
	r.Post("/editpost/:userid/:postid", controller.EditPost)
	r.Post("/deletepost/:userid/:postid", controller.DeletePost)
	r.Get("/getpost/:userid/:postid", controller.GetPosts)

	r.Post("/addcomment/:userid/:postid", controller.AddCommentToPost)
	r.Post("/addsubcomment/:userid/:postid/:commentid", controller.AddCommentToComment)

	r.Post("/likepost/:userid/:postid", controller.LikeToPost)
	r.Post("/likecomment/:userid/:postid/:commentid", controller.LikeToComment)

	r.Post("/addfriend/:userid/:friendid", controller.AddFriend)
	r.Post("/removefriend/:userid/:friendid", controller.RemoveFriend)

	r.Post("/createshorturl", controller.CreateShortURL)
	r.Post("/redirect/:shortkey", controller.RedirectShortURL)

	r.Listen(":3000")
}
