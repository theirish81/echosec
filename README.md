# EchoSec
A Golang middleware for the [Labstack Echo Server](https://echo.labstack.com/) that simplifies the process of describing
the prerequisites for a request to access combinations of endpoints and methods.
The most common application is to offload boring and copy-paste security constraints to a middleware.

**Ie.** *Can a user with a given JWT token perform `DELETE /user/:id`?*

## Example

**Premises:**

* `GetClaims(c)` is a dummy function that represents, in the example, the retrieval of JWT claims.
* 

```go
import "github.com/theirish/echosec"

/*...*/

m := echosec.Middleware(echosec.Config{
    PathMapping: echosec.PathItems{
        {
            Patterns: echosec.Patterns{"/api/v1/user"},
            PathValidation: func(c echo.Context) error {
                if GetClaims(c).Admin {
                    return nil
                }
                return errors.NewForbiddenError()
            },
        },
        {
            Patterns: echosec.Patterns{"/api/v1/user/:userId"},
            PathValidation: func(c echo.Context) error {
                if GetClaims(c).CanAdminUserData(c.Param("userId")) {
                    return nil
                }
                return errors.NewForbiddenError()
            },
        },
        {
            Patterns: echosec.Patterns{"/api/v1/workspace"},
            Methods: echosec.ValidationMap{
                "GET": func(c echo.Context) error {
                    if GetClaims(c).Admin {
                        return nil
                    }
                    return errors.NewForbiddenError()
                },
            },
            PathValidation: func(c echo.Context) error {
                return nil
            },
        },
        {
            Patterns: echosec.Patterns{"/api/v1/workspace/:workspaceId",
                "/api/v1/workspace/:workspaceId/membership",
                "/api/v1/workspace/:workspaceId/membership/:membershipId"},
            Methods: echosec.ValidationMap{
                "GET": func(c echo.Context) error {
                    if GetClaims(c).CanAccessWorkspace(c.Param("workspaceId")) {
                        return nil
                    }
                    return errors.NewForbiddenError()
                },
                "PUT,DELETE": func(c echo.Context) error {
                    if GetClaims(c).CanAdminWorkspace(c.Param("workspaceId")) {
                        return nil
                    }
                    return errors.NewForbiddenError()
                },
            },
        },
    },
})

/*...*/

e := echo.New()
e.Use(m)
```