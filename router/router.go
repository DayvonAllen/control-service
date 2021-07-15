package router

import (
	"example.com/app/handlers"
	"example.com/app/middleware"
	"example.com/app/repo"
	"example.com/app/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRoutes(app *fiber.App) {
	ch := handlers.CommentHandler{CommentService: services.NewCommentService(repo.NewCommentRepoImpl())}
	sh := handlers.StoryHandler{StoryService: services.NewStoryService(repo.NewStoryRepoImpl())}
	reh := handlers.ReplyHandler{ReplyService: services.NewReplyService(repo.NewReplyRepoImpl())}
	uh := handlers.UserHandler{UserService: services.NewUserService(repo.NewUserRepoImpl())}
	ah := handlers.AuthHandler{AuthService: services.NewAuthService(repo.NewAuthRepoImpl())}

	app.Use(recover.New())
	api := app.Group("", logger.New())

	stories := api.Group("application/storage/app/stories")
	stories.Get("/:id", middleware.IsLoggedIn, sh.FindStory)
	stories.Delete("/:id", middleware.IsLoggedIn, sh.DeleteStory)
	stories.Get("/", middleware.IsLoggedIn, sh.FindAll)

	comments := api.Group("application/storage/app/comment")
	comments.Delete("/:id",middleware.IsLoggedIn, ch.DeleteById)

	reply := api.Group("application/storage/app/reply")
	reply.Delete("/:id", middleware.IsLoggedIn, reh.DeleteById)

	auth := api.Group("application/storage/app/auth")
	auth.Post("/login", ah.Login)

	user := api.Group("application/storage/app/users")
	user.Get("/", middleware.IsLoggedIn, uh.GetAllUsers)
	user.Delete("/delete",middleware.IsLoggedIn,  uh.DeleteByID)
}

func Setup() *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		ExposeHeaders: "Authorization",
	}))

	SetupRoutes(app)
	return app
}
