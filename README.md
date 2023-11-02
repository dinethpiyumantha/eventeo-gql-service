## "EVENTEO" Event management service

###### Create

```
mutation CreateEventListing($input: CreateEventListingInput!) {
  createEventListing(input:$input){
    _id
    title
    description
    organizer
    url
  }
}

# variables
{
  "input": {
    "title": "Software development competition",
    "description": "This is a sample event",
    "organizer": "dineOrg",
    "url": "dinethpiyumantha.netlify.app"
  }
}
```

###### Get a event

```
query GetEvent($id: ID!){
  event(id:$id){
    _id
    title
    description
    url
    organizer
  }
}

# variables
{
  "id": "654369d97c80cd8aebdf77a4"
}
```

###### Update

```
mutation UpdateEvent($id: ID!, $input: UpdateEventListingInput!) {
  updateEventListing(id:$id, input:$input){
    title
    description
    _id
    organizer
    url
  }
}

#variables
{
  "id": "654369d97c80cd8aebdf77a4",
  "input": {
    "title": "Software development competition 2024"
  }
}
```

###### Delete

```
mutation DeleteQuery($id: ID!){
  deleteEventListing(id:$id){
    deleteEventId
  }
}

#variables
{
  "id": "654369d97c80cd8aebdf77a4"
}
```
