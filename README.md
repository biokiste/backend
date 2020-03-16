# backend
✨New website backend of Biokiste e.V.✨

## development

expects connection string of mysql instance in config.toml (app root)

- compile backend with `go build`
- run backend with `./backend`

## new database tables

### Settings

```sql
CREATE TABLE `Settings` (
  `ID` int(11) NOT NULL,
  `ItemKey` varchar(255) NOT NULL,
  `ItemValue` varchar(255) NOT NULL,
  `CreatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreatedBy` int(11) NOT NULL,
  `UpdatedAt` datetime DEFAULT NULL,
  `UpdatedBy` int(11) DEFAULT NULL,
  `UpdateComment` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

### api 

Create user:
- POST to `/api/user/auth/create` with body:
`{
	"email":      "tina@teewurst.org",
	"password": "**********",
	"lastname": "teewurst",
	"firstname":   "tina",
	"mobile": "8348349",
	"street": "Fleischergasse",
	"credit_date": "2018-03-12"
}`


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

Update user:
- PATCH to `/api/user` with body:
`{	
	"id": 1,
	"username": "ro.ri",
	"email": "roland.rindfleisch@web.de",
	"lastname": "Rindfleisch",
	"firstname": "Roland",
	"mobile": "2837432847",
	"street": "Industriestraße 101",
	"zip": "04229",
	"city": "Leipzig",
	"date_of_birth": "1901-08-19",
	"date_of_entry": "2020-03-03"
}`

for other routes see @ `routes.go`
