
```bash
git clone https://github.com/hariadivicky/nano-example
```

Now we assume that you are on `nano-example/` directory.

Then you can create new SQLite database on `database/`

```bash
touch database/todos.db
```

Feel free to choose your database name. You could find this configuration on `main.go` file

```go
models.OpenDatabase("file:./database/todo.db")
```
```bash
$  sqlite3 database/todo.db
>  .read database/create_todos_table.sql
>  .exit
```
```bash
go build && ./todo
```

server runs on `http://localhost:8080`


show the collection of todo.

`GET` [http://localhost:8080/serviceregistry/echo]

response:

```json
{
    "collection": [
        {
            "id": 1,
            "title": "learn golang",
            "is_done": false
        },
        {
            "id": 2,
            "title": "practice golang",
            "is_done": false
        }
    ]
}
```

### Create Todo

create new todo

`POST` [http://localhost:8080/serviceregistry/query]

Headers:

`Content-Type: application/x-www-form-urlencoded`

`Data: title=use%20nano&is_done=true`

response:

```json
{
    "todo": {
        "id": 3,
        "title": "use nano",
        "is_done": true
    }
}
```

### Todo Detail

display todo detail

`POST` [http://localhost:8080/serviceregistry/register]

response:

Headers:

`Content-Type: application/x-www-form-urlencoded`

`Data: title=use%20nano&is_done=true`

response:

```json
{
    "todo": {
        "id": 3,
        "title": "use nano",
        "is_done": true
    }
}
```



### Remove Todo

delete todo record

`DELETE` [http://localhost:8080/serviceregistry/unregister]

response:

```json
{
    "deleted": true
}
```
