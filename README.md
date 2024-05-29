# EchoSec
A Golang middleware for the [Labstack Echo Server](https://echo.labstack.com/) that simplifies the process of describing
the prerequisites for a request to access combinations of endpoints and methods, using an OpenAPI extension to direct
the routing and Go code to determine the rules. The most common application is to offload boring and copy-paste
security constraints to a middleware.

**I.e.** *Can a user with a given JWT token perform `DELETE /user/:id`?*

## The OpenAPI specification
The first thing to do is to annotate an OpenAPI file with the `x-echosec` custom tag. Each operation requiring the
attention of EchoSec needs to be annotated.

**I.e.** 
```yaml
"/user/{userId}":
  get:
    summary: |-
        returns the user details. use `/me` to get current user. Users can only fetch their own details, admins, instead,
        can see everything.
    operationId: getUser
    x-echosec:
        function: can_admin_user
```
The key `x-echosec.function` defines an arbitrary ID that will be mapped to a function. The string can be anything,
but we advise to take the time to define all the various access modes first, and then diligently assign them.
In this case, `can_admin_user` represents the need of validating the requesting user and making sure she's allowed
to manage the requested user.

Additionally, you can specify further params. Params are just labels that will be passed to the function, to distinguish
sub-behaviors.

**I.e.**
```yaml
"/workspace/{workspaceId}/_compute-organizations":
    post:
        operationId: computeWorkspaceOrganizations
        x-echosec:
          function: workspace_user_read
          params:
            - premium
```
In this example, the main behavior is `workspace_user_read`, however, we further specify the `premium` param so that
the function knows there's a sub-behavior (premium) to be taken into account.

## The Go Code
We first set up our configuration as in:
```go
cfg, _ := echosec.NewOApiConfig(openApiBytes, map[string]echosec.OApiValidationFunc{
    // all the functions mentioned in this snippet, such as 'getClaims' are fictitious and presented to 
    // offer an example that you can relate to.
    "can_admin_user": func(c echo.Context, params []string) error {
        if GetClaims(c).CanAdminUserData(c.Param("userId")) {
            return nil
        }
        return errors.NewForbiddenError()
    },
    "workspace_user_read": func(ctx echo.Context, params []string) error {
        claims := GetClaims(ctx)
        if claims.Admin {
            return nil
        }
        if slices.Contains(params, "premium") && !claims.IsPremiumWorkspace(ctx.Param("workspaceId")) {
            return errors.NewForbiddenError()
        }
        if claims.CanAccessWorkspace(ctx.Param("workspaceId")) {
            return nil
        }
        return errors.NewForbiddenError()
    },
}
```
So as you can se we're mapping the OpenAPI `x-echosec.function` items to actual go functions. If the user is allowed
to perform a certain action, we return `nil`. If the user is not, then we return the appropriate error.
The `openApiBytes` argument is an serialized OpenAPI file. It can either be plain text or GZipped.

## Wiring all together
It's pretty simple as it works like any other Echo Framework middlewares.
```go
    server := echo.New()
    server.Use(echosec.WithOpenApiConfig(cfg))
```