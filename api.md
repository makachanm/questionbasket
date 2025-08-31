# API Design

## Design

All POST endpoint uses JSON object to transfer and serve data.
Invalid JSON Object can cause any kind of error of malfunction of API.

## Types & Constants

- QID
 uses string, use for identify question.

- AID
 uses string, use for one-time token for answer the question.

- Share Range (for Fediverse option)
```
1 = Public Share 
2 = Home Share
3 = Private Share
```

## Question

### Add new questions
```
 POST /question/make
```

add new question to specific user.

***Needs***

`content` uses string
`is_nsfw` uses boolean


***Params***
```
{
   "content": "foobar",
   "is_nsfw": false
}
```

***Returns***
```
{
   "qid": (QID)
}
```

<hr />

### Reply a answer to questions
```
 POST /question/answer?aid={AID}
```

reply a answer to specific question.

***Needs***

`id`
uses AID

`content`
uses string

`share_range`
uses int

***Params***
```
{
   "content": "foobar",
   "share_range": 1
}
```

***Returns***
```
{
   "qid": (QID)
}
```


<hr />

### Get Answered Questions
```
 GET /qustion/{QID}
```
return recent 20 questions.

***Optional***

`id`
uses QID

***Returns***
```
[
   {
      "content": "foobar",
      "created_at": XXXXXXXX,
      "answer": "fizzbuzz",
      "share_range": 1,
      "id": (QID),
   },
   (....)
]
```

<hr />

### Get Profile
```
 GET /profile
```


***Returns***
```
{
   "name": "foobar",
   "description": "fizzbuzz",
   "based_on": 1 
}
```

