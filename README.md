# go-okta
incomplete okta client written in golang


# Usage

    client := NewClient("api token", "dev-532085", true, http.DefaultClient)
    ctx := context.Background()
    ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
    defer cancel()
    group, err := client.Group(ctx, "groupId")
