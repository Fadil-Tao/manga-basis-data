# Manga Basis Data

### Note 
The project is not complete, todo :
- add migrations
- default data
- dockerize
- simple frontend with templating

### Description

This is Rest A.P.I project written in go and mysql with pure net/http. This project logic handled with mysql stored procedure.

### Prequisites

1. Go Installed 
2. Mysql Installed

### Installation

1. Git Clone `https://github.com/Fadil-Tao/manga-basis-data.git`
2. Create a user in mysql 
3. Follow the env example
4. run api with `go run cmd/api/main.go`

--------------------

### Auth

Authentication-related endpoints for login, registration, and logout.

<details>
<summary><code>POST</code> <code><b>/login</b></code></summary>

**Description:** User login.

**Example Request:**

```json
{
  "email": "rio@gmail.com",
  "password": "rionandosoeksin"
}
```

**Example Response**

```json
{
  "Message": "Login Success"
}
```

</details>

<details>
<summary><code>POST</code> <code><b>/register</b></code></summary>

**Description:** User registration.

**Example Request** example

```json
{
  "username": "ilham",
  "email": "ilham@gmail.com",
  "password": "ilham123456"
}
```

**Example Response** example

```json
{
  "message": "User Registered successfully"
}
```

</details>
<details>
<summary><code>POST</code> <code><b>/logout</b></code></summary>

**Description:** User logout.

**Example Response**

```json
{
  "Message": "Log out Success"
}
```
</details>

----------------------------------------

### Author

Endpoints for managing author data, including adding, retrieving, updating, and deleting authors.

<details>
<summary><code>POST</code> <code><b>/author</b></code></summary>

**Description:** Add a new author (Admin only).

**Example Request:**

```json
{
  "name": "miyamoto",
  "birthday": "1960-06-07",
  "biography": "Guy who likes to draw"
}
```

**Example Response**

```json
{
  "message": "Author Created successfully"
}
```

</details>

<details>
<summary><code>GET</code> <code><b>/author</b></code></summary>

**Description:** Retrieve all authors

**Example Response :**

```json
{
  "message": "Authors Succesfully Retrieved",
  "data": [
    {
      "id": "1",
      "name": "Kinji Hakari",
      "birthday": "1982-10-02"
    },
    {
      "id": "3",
      "name": "fujimoto",
      "birthday": "1998-05-20"
    }
  ]
}
```

</details>

<details>
<summary><code>GET</code> <code><b>/author/{id}</b></code></summary>

**Description :** Retrieve author details by ID.

**Example Response :**

```json
{
  "message": "data succesfully retrieved",
  "data": {
    "id": "10",
    "name": "hiroiko tanaka",
    "birthday": "1998-01-20",
    "biography": "i like fish",
    "Manga": null
  }
}
```

</details>

<details> <summary><code>PUT</code> <code><b>/author/{id}</b></code></summary>

**Description :** Update author details (admin only).

**Example Request :**

```json
{
  "name": "hiroiko tanaka",
  "birthday": "1998-01-20",
  "biography": "i like fish"
}
```

**Example Response :**

```json
{
  "name": "hiroiko tanaka",
  "birthday": "1998-01-20",
  "biography": "i like fish"
}
```

</details>

<details>
<summary><code>DELETE</code> <code><b>/author/{id}</b></code></summary>

**Description :** Delete an author by ID (admin only).
**Example Response:**

```json
{
  "message": "Author successfully deleted"
}
```

</details>

------------------------------

### Genre

Endpoints for managing genre data, including creating, retrieving, updating, and deleting genres.

<details>
<summary><code>POST</code> <code><b>/genre</b></code></summary>

**Description:** Create a new genre (Admin Only)

**Request:**

```json
{
  "name": "romance",
  "description": "Story about loves"
}
```

**Response:**

```json
{
  "message": "Genre Created successfully"
}
```

</details>

<details>
<summary><code>GET</code> <code><b>/genre</b></code></summary>

**Description :** Retrieve All Genre
**Example Response :**

```json
{
  "message": "Genres successfully retrieved",
  "data": [
    {
      "id": "1",
      "name": "horror",
      "description": "stories about scary stuffs"
    },
    {
      "id": "4",
      "name": "sport",
      "description": "about sports"
    },
    {
      "id": "6",
      "name": "Isekai",
      "description": "Story about transport into different world"
    },
    {
      "id": "8",
      "name": "motor jawa",
      "description": "Story about transport into different world"
    },
    {
      "id": "9",
      "name": "romance",
      "description": "Story about loves"
    }
  ]
}
```

</details>

<details> <summary><code>POST</code> <code><b>/genre/{id}</b></code></summary>

**Description :** Update genre data by ID (Admin Only)
**Example Request:**

```json
{
  "name": "native fantasy",
  "description": "stories about native"
}
```

**Example Response:**

```json
{
  "message": "Genre updated successfully"
}
```

</details>

<details>
<summary><code>DELETE</code> <code><b>/genre/{id}</b></code></summary>

**Description :** Delete Genre By Id (Admin only)
**Example Response :**

```json
{
  "message": "Genre successfully deleted"
}
```
</details>

----------------------------

### Manga

Endpoints for managing manga, including creation, association with authors/genres, liking, retrieval, updating, and deletion.

<details>
<summary><code>POST</code> <code><b>/manga</b></code></summary>

**Description:** Create a new manga (Admin Only)

**Example Request:**

```json
{
  "title": "Berserkin time",
  "synopsys": "This one is so edgy",
  "status": "in_progress",
  "published_at": "2004-05-22",
  "finished_at": "2016-09-22"
}
```

**Example Response:**

```json
{
  "message": "Manga Created successfully"
}
```

</details>
 <details> <summary><code>POST</code> <code><b>/manga/{id}/author</b></code></summary>

**Description :** Associate an author with a manga (Admin Only)
**Example Request :**

```json
{
  "authorId": 5
}
```

**Example Response :**

```json
{
  "message": "Author and Manga connected successfully"
}
```

</details>

<details>
 <summary><code>POST</code> <code><b>/manga/{id}/genre</b></code></summary>

**Description :** Associate a genre with a manga
**Example Request :**

```json
{
  "genreId": 1
}
```

**Example Response:**

```json
{
  "message": "Toggle triggered successfully"
}
```

</details>

<details>
   <summary><code>GET</code> <code><b>/manga</b></code></summary>

**Description :** Retrieve a list of all manga.
**Example Reponse :**

```json
{
  "message": "succefully retrieved manga",
  "data": [
    {
      "id": "2",
      "title": "Attack On Titan",
      "status": "finished",
      "published_at": "2010-05-22",
      "finished_at": "2020-09-22",
      "rating": 0,
      "totalReview": 0,
      "likes": 1,
      "totalUserRated": 0
    },
    {
      "id": "4",
      "title": "Berserk",
      "status": "finished",
      "published_at": "2004-05-22",
      "finished_at": "2016-09-22",
      "rating": 0,
      "totalReview": 0,
      "likes": 1,
      "totalUserRated": 0
    }
  ]
}
```

</details>

<details> 
<summary><code>GET</code> <code><b>/manga?name={name}</b></code></summary>
 
**Description :** Search for manga by name.

**Example:** `/manga?name=ber`

**Example Response :**

```json
{
  "message": "succefully retrieved manga",
  "data": [
    {
      "id": "4",
      "title": "Berserk",
      "status": "finished",
      "published_at": "2004-05-22",
      "finished_at": "2016-09-22",
      "rating": 0,
      "totalReview": 0,
      "likes": 0,
      "totalUserRated": 0
    }
  ]
}
```

</details>

<details>
    <summary><code>GET</code> <code><b>/manga/{id}</b></code></summary>

**Description:** Retrieve manga details by ID.

**Example Response :**

```json
{
  "message": "manga success retrieved",
  "data": {
    "id": "4",
    "title": "Berserk",
    "synopsys": "This one is so edgy",
    "status": "finished",
    "published_at": "2004-05-22",
    "finished_at": "2016-09-22",
    "genre": [
      {
        "id": "1",
        "name": "horror"
      }
    ],
    "author": [
      {
        "id": "5",
        "name": "Hajime Isayama"
      }
    ]
  }
}
```

</details>

<details>
<summary><code>PUT</code> <code><b>/manga/{id}</b></code></summary>

**Description :** Retrieve author details by ID.

**Example:** `/manga/4`

**Example Request:**

```json
{
  "title": "ben ten",
  "synopsys": "naruto",
  "status": "in_progress",
  "published_at": "2000-05-22",
  "finished_at": "2015-09-22"
}
```

**Example Response:**

```json
{
  "message": "Manga updated successfully"
}
```

</details>

<details>
<summary><code>DELETE</code> <code><b>/manga/{id}</b></code></summary>

**Description :** Retrieve author details by ID.
**Example :** `/manga/5`

**Example Response :**

```json
{
  "message": "Manga successfully deleted"
}
```

</details>

<details>
<summary><code>DELETE</code> <code><b>/manga/{mangaId}/author/{authorId}</b></code></summary>

**Description :** Delete an association between a manga and an author. (Admin)

**Example :** `/manga/1/author/7`

**Example Response :**

```json
{
  "message": "Deleted successfully"
}
```

</details>

<details>
<summary><code>DELETE</code> <code><b>/manga/{mangaId}/genre/{genreId}</b></code></summary>

**Description :** Delete an association between a manga and a genre.(Admin Only)

**Example Response :**

```json
{
  "message": "Deleted successfully"
}
```

</details>

### Rating

Endpoint for user to give rating to a manga on scale 1-10

<details>
<summary><code>POST</code> <code><b>/manga/{id}/rating</b></code></summary>

**Description:** Rate a manga.

**Example:** `/manga/1/rating`

**Example Request:**

```json
{
  "rating": 10
}
```

</details>

--------------------------

### Readlist

Endpoints for managing readlists and their associated manga items.

<details>
<summary><code>POST</code> <code><b>/readlist</b></code></summary>

**Description:** Create a new readlist.

**Example Request:**

```json
{
  "name": "xxx",
  "description": "xxx"
}
```

**Example Response :**

```json
{
  "message": "Readlist Created successfully"
}
```

</details>
<details> 
<summary><code>POST</code> <code><b>/readlist/{id}/item</b></code></summary>

**Description :** : Add a manga to a readlist (Admin Only)
**Example Request:**

```json
{
  "mangaId": "2",
  "readStatus": "done"
}
```

**Example Response:**

```json
{
  "message": "manga added to readlist successfully"
}
```

</details>

<details>
<summary><code>GET</code> <code><b>/readlist</b></code></summary>

**Description :** Retrieve all readlists

**Example Response:**

```json
{
  "message": "readlist succesfully retrieved",
  "data": [
    {
      "id": "1",
      "owner": "jenipers",
      "name": "My most despised readlist",
      "description": "i hate this book very much",
      "created_at": "2024-11-29 05:24:44",
      "updated_at": "2024-11-29 05:54:30"
    },
    {
      "id": "3",
      "owner": "rionandoo",
      "name": "isekaioo",
      "description": "this is just bunch of book i wish to read if i have free time",
      "created_at": "2024-12-03 03:49:25",
      "updated_at": "2024-12-03 03:49:25"
    }
  ]
}
```

</details>

<details>
<summary><code>GET</code> <code><b>/readlist/{id}/item</b></code></summary>

**Description :** Retrieve all manga items in a readlist.

**Example Response :**

```json
{
  "message": "readlist item succesfully retrieved",
  "data": [
    {
      "mangaId": "4",
      "title": "Berserk",
      "readStatus": "done",
      "addedAt": "2024-12-03 03:52:08"
    }
  ]
}
```

</details>

<details>
<summary><code>PUT</code> <code><b>/readlist/{id}</b></code></summary>

**Description :** Update a readlist's details

**Example Request :**

```json
{
  "name": "Koleksi Buku Rio",
  "description": "People Come and go"
}
```

**Example Response :**

```json
{
  "message": "Readlist updated successfully"
}
```

</details>

<details>
<summary><code>PUT</code> <code><b>/readlist/{readlistId}/manga/{mangaId}</b></code></summary>

**Description :** Update the reading status of a manga in a readlist.

**Example Request**

```json
{
  "status": "done"
}
```

**Example Response**

```json
{
  "message": "Readlist updated successfully"
}
```

</details>

<details>
<summary><code>DELETE</code> <code><b>/readlist/{readlistId}/manga/{mangaId}</b></code></summary>

**Example Response:**

```json
{
  "message": "Readlist item deleted successfully"
}
```

</details>

--------------------------

### Manga Reviews

Endpoints for creating, retrieving, and managing reviews for manga.

<details>
<summary><code>POST</code> <code><b>/manga/{id}/review</b></code></summary>

**Description:** Create a review for a manga.

**Example:** `/manga/4/review`

**Example Request:**

```json
{
  "review": "My fav so far",
  "tag": "Reccomended"
}
```

**Example Response :**

```json
{
  "message": "Review Created successfully"
}
```

</details>

<details>
<summary><code>GET</code> <code><b>/manga/{id}/review</b></code></summary>

**Description :** Retrieve a list of reviews for a manga.
**Example :** `/manga/4/review`
**Example Response :**

```json
{
  "message": "review succesfully retrieved",
  "data": [
    {
      "username": "fery",
      "user_id": "2",
      "review": "i dont love it!",
      "tag": "Not Reccomended",
      "created_at": "2024-12-03 07:29:49",
      "like": 0
    },
    {
      "username": "ilham",
      "user_id": "5",
      "review": "i love it!",
      "tag": "Reccomended",
      "created_at": "2024-12-03 07:29:18",
      "like": 0
    },
    {
      "username": "rionandoo",
      "user_id": "4",
      "review": "i hate narutoo",
      "tag": "Mixed Feelings",
      "created_at": "2024-12-03 05:39:37",
      "like": 1
    }
  ]
}
```

</details>

<details>
<summary><code>GET</code> <code><b>/manga/{mangaId}/review/{reviewerId}</b></code></summary>

**Description :** Retrieve a specific review for a manga by the reviewer's ID (Admin)

**Example :** `/manga/4/review/4`

```json
{
  "message": "review data retrieved succesfully",
  "data": {
    "manga_id": "4",
    "username": "rionandoo",
    "user_id": "4",
    "review": "i hate narutoo",
    "tag": "Mixed Feelings",
    "created_at": "2024-12-03 05:39:37",
    "like": 1
  }
}
```

</details>

<details>
<summary><code>PUT</code> <code><b>/manga/{mangaId}/review/{reviewerId}/like</b></code></summary>

**Description :** Update a review for a manga (Admin Only)

**Example :** `/manga/2/review/4`

**Example Request :**

```json
{
  "review": "i hate narutoo",
  "tag": "Not Reccomended"
}
```

**Example Response :**

```json
{
  "message": "Review Updated successfully"
}
```

</details>

----------------------------

### User

Endpoints for managing user accounts, retrieving user-related information, and performing user-related operations

<details>
<summary><code>GET</code> <code><b>/users</b></code></summary>

**Description :** Retrieve all users.

**Example Response :**

```json
{
  "message": "user succesfully retrieved",
  "data": [
    {
      "id": 2,
      "username": "fery",
      "email": "fery@gmail.com",
      "created_at": "2024-11-13 21:31:34"
    },
    {
      "id": 3,
      "username": "jenipers",
      "email": "jenipers@yahoo.com",
      "created_at": "2024-11-25 19:59:36"
    },
    {
      "id": 4,
      "username": "rionandoo",
      "email": "rio@gmail.com",
      "created_at": "2024-12-02 14:20:07"
    },
    {
      "id": 5,
      "username": "ilham",
      "email": "ilham@gmail.com",
      "created_at": "2024-12-03 01:38:48"
    }
  ]
}
```

</details>

<details>
<summary><code>GET</code> <code><b>/users?username={query}</b></code></summary>

**Description :** Search for a user by username
**Example :** `/users?username=rio`

**Example Response :**

```json
{
  "message": "user succesfully retrieved",
  "data": [
    {
      "id": 4,
      "username": "rionandoo",
      "email": "rio@gmail.com",
      "created_at": "2024-12-02 14:20:07"
    }
  ]
}
```

</details>

<details>
<summary><code>GET</code> <code><b>/users/{username}</b></code></summary>

**Description :** Retrieve detailed information about a specific user.

**Example :** `/users/fery`

**Example Response :**

```json
{
  "message": "user succesfully retrieved",
  "data": {
    "id": 2,
    "username": "fery",
    "email": "fery@gmail.com",
    "created_at": "2024-11-13 21:31:34"
  }
}
```

</details>

<details>
<summary><code>GET</code> <code><b>/users/{username}/likedmanga</b></code></summary>

**Description :** Retrieve the list of manga liked by a specific user.

**Example :** `/users/rionandoo/likedmanga`

**Example response :**

```json
{
  "message": "mangas succesfully retrieved",
  "data": [
    {
      "id": "4",
      "title": "Berserk",
      "status": "finished",
      "published_at": "2004-05-22",
      "finished_at": "2016-09-22",
      "likedAt": "2024-12-03 03:25:16",
      "likes": 1
    }
  ]
}
```

</details>

<details>
<summary><code>GET</code> <code><b>/users/{username}/ratedmanga</b></code></summary>

**Description :** Retrieve the list of manga rated by a specific user.

**Example :** `/users/rionandoo/ratedmanga`

**Example Response :**

```json
{
  "message": "mangas succesfully retrieved",
  "data": [
    {
      "id": "4",
      "title": "Berserk",
      "status": "finished",
      "published_at": "2004-05-22",
      "finished_at": "2016-09-22",
      "ratededAt": "2024-12-03 03:45:35",
      "yourRating": 10,
      "rating": 1,
      "totalUserRated": 10
    },
    {
      "id": "1",
      "title": "naruto",
      "status": "in_progress",
      "published_at": "2000-05-22",
      "finished_at": "2015-09-22",
      "ratededAt": "2024-12-02 22:35:57",
      "yourRating": 7,
      "rating": 1,
      "totalUserRated": 7
    }
  ]
}
```

</details>

<details>
<summary><code>GET</code> <code><b>/users/{username}/readlist</b></code></summary>

**Description :** Retrieve the list of readlists owned by a user.

**Example :** `/users/rionandoo/readlist`

**Example Response :**

```json
{
  "message": "mangas succesfully retrieved",
  "data": [
    {
      "id": "3",
      "name": "Koleksi Buku Rio",
      "description": "People Come and go",
      "created_at": "2024-12-03 03:49:25",
      "updated_at": "2024-12-03 03:59:51"
    },
    {
      "id": "5",
      "name": "yohsa",
      "description": "this is just bunch of book i wish to read if i have free time",
      "created_at": "2024-12-03 05:09:52",
      "updated_at": "2024-12-03 05:09:52"
    }
  ]
}
```
</details>

<details>
<summary><code>PUT</code> <code><b>/users/{username}</b></code></summary>

**Description :** Update a user's information.

**Example :** `/users/rionandoo`

**Example request :**
```json
{
    "username": "rionandoo"
}
```

**Example response :**
```json
{
    "message": "User Updated successfully"
}
```
</details>

<details>
<summary><code>DELETE</code> <code><b>/users/{username}</b></code></summary>

**Description :** Delete a user
**Example :** `/users/rionandoo`

**Example Response:**
```json
{
    "message": "user deleted Successfully"
}
```
</details>