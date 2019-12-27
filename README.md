# backend
✨New website backend of Biokiste e.V.✨

## development

expects running mysql instance on `localhost:8889` (root, root) with db foodkoop_biokiste

- compile backend with `go build`
- run backend with `./backend`


### api 

Insert user transactions:
- POST to `/api/transaction` with body:
`{
	"transactions": [
	{	
		"amount": 100.00,
		"created_at": "2019-12-27 17:30",
		"category_id": 1,
		"status": 1
	}
	],
	"user": {
	"id": 176
	}
}`

Update doorcode:
- PATCH to `/api/settings/doorcode` with body:
`{
	"doorcode": "Außen: 225588 Innen:685259",
	"updated_at": "2019-12-23 14:00",
	"updated_by": 176
}`

for other routes see @ `routes.go`