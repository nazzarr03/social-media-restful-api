package controller

import (
	"github.com/gofiber/fiber"
	"github.com/nazzarr03/social-media-restful-api/database"
	"github.com/nazzarr03/social-media-restful-api/models"
)

func LikeToPost(c *fiber.Ctx) {
    user := models.User{}
    post := models.Post{}

    database.Db.First(&user, "user_id = ?", c.Params("userid"))

    if user.UserID == 0 {
        c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "User not found",
        })
        return
    }

    database.Db.First(&post, "post_id = ?", c.Params("postid"))

    if post.PostID == 0 {
        c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Post not found",
        })
        return
    }

    like := models.Like{
        UserID: user.UserID,
        PostID: post.PostID,
        IsLiked: true,
    }

    database.Db.Create(&like)

	if database.Db.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create like to post",
		})
		return
	}

    c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Post liked",
    })
}

func LikeToComment(c *fiber.Ctx) {
    user := models.User{}
    post := models.Post{}
    comment := models.Comment{}

    database.Db.First(&user, "user_id = ?", c.Params("userid"))

    if user.UserID == 0 {
        c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "User not found",
        })
    }

    database.Db.First(&post, "post_id = ?", c.Params("postid"))

    if post.PostID == 0 {
        c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Post not found",
        })
    }
    
    database.Db.First(&comment, "comment_id = ?", c.Params("commentid"))

    if comment.CommentID == 0 {
        c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Comment not found",
        })
    }

    like := models.Like{
        UserID: user.UserID,
        PostID: post.PostID,
        CommentID: &comment.CommentID,
        IsLiked: true,
    }

    database.Db.Create(&like)

    if database.Db.Error != nil {
        c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Cannot create like to comment",
        })
    }

    c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Comment liked",
    })
}