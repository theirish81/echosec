# EchoSec

**EchoSec** is a Golang middleware for the [Labstack Echo Server](https://echo.labstack.com/) that simplifies the process of enforcing request prerequisites across endpoint-method combinations. It uses an OpenAPI extension to direct routing and custom Go functions to define access rules.

Its primary use case is offloading repetitive, boilerplate security logic to middleware. On top of this, EchoSec introduces an advanced *request labeling* system, allowing you to define rules that dynamically assign labels to a request context when matched.

> **Example:** *Can a user with a given JWT token perform `DELETE /user/:id`?*

---

## OpenAPI Integration

Start by annotating your OpenAPI file with the custom `x-echosec` tag. Any operation requiring EchoSec enforcement must include this annotation.

### Core: Functions

```yaml
"/user/{userId}":
  get:
    summary: |-
        Returns the user details. Use `/me` to get the current user. Users can only fetch their own details; admins can see everything.
    operationId: getUser
    x-echosec:
      function: can_admin_user
```

The `x-echosec.function` key defines an identifier that maps to a Go function. While the string is arbitrary, we recommend defining all access modes up front and assigning them consistently.

In this example, `can_admin_user` refers to a function that verifies the requesting user's ability to manage the targeted user.

---

### Parameters

You can further specify behavior using `params`. These are simple strings passed to the validation function to distinguish sub-cases.

```yaml
"/workspace/{workspaceId}/_compute-organizations":
  post:
    operationId: computeWorkspaceOrganizations
    x-echosec:
      function: workspace_user_read
      params:
        - premium
```

Here, the base function `workspace_user_read` applies, but the `premium` param signals a more specific behavior to enforce.

---

### Labels

You can assign labels to the request context using conditional expressions:

```yaml
"/workspace/{workspaceId}/_compute-organizations":
  post:
    operationId: computeWorkspaceOrganizations
    x-echosec:
      function: workspace_user_read
      params:
        - premium
      labels:
        - label: "agent"
          condition: '"agent" in claims.Roles'
```

The condition syntax follows the [Expr](https://github.com/expr-lang/expr) language.

---

## Go Code

### Function Mapping

First, configure your validation logic:

```go
cfg, _ := echosec.NewOApiConfig(openApiBytes, map[string]echosec.OApiValidationFunc{
    // Example functions—adapt to your actual claims and access model.
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
}, true)
```

Here, the `x-echosec.function` values in the OpenAPI spec are mapped to actual Go functions. These should return `nil` if access is granted, or an error if denied.

* `openApiBytes` is your OpenAPI spec in serialized form (either plain text or gzipped).
* The final `true` argument enables **request validation**. If enabled, malformed requests will be rejected before reaching your access rules.

---

### Labels in Go

Label conditions rely on data extracted from the Echo context. You can promote certain values to the top-level expression scope with `.WithVars()`:

```go
cfg.WithVars("Claims")
```

This allows your label conditions to access `Claims` directly.

The evaluated labels are stored in the `EchoSecContext`, which is made available via:

```go
ectx.Get("echosecContext")
// or
ctx.Value("echosecContext")
```

The `EchoSecContext` struct looks like this:

```go
type EchoSecContext struct {
    Config OApiEchoSec
    Labels []string
}
```

* `Config` refers to the matched security configuration.
* `Labels` contains the list of active labels assigned to the request.

---

## Response Validation (Optional)

EchoSec also supports validating responses against the OpenAPI spec. This ensures your application complies with the expected schema.

> **Note:** Request validation must be enabled for this to work.

To validate a response, use `SecBlob` instead of `ctx.JSON` or `ctx.Blob`:

```go
echosec.SecBlob(ctx, http.StatusOK, map[string]any{
    "a": "foobar",
    "b": 1,
})
```

* The first parameter is the Echo context.
* The second is the HTTP status code.
* The third is the response payload (map, string, or `[]byte`).

If the payload doesn’t match the OpenAPI schema, EchoSec returns an error.

---

## Putting It All Together

Using EchoSec is as simple as plugging it in as a middleware:

```go
server := echo.New()
server.Use(echosec.WithOpenApiConfig(cfg))
```

---