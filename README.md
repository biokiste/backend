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
  `SettingKey` varchar(255) NOT NULL,
  `Value` varchar(255) NOT NULL,
  `Description` varchar(255) NOT NULL,
  `CreatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreatedBy` int(11) NOT NULL,
  `UpdatedAt` datetime DEFAULT NULL,
  `UpdatedBy` int(11) DEFAULT NULL,
  `UpdateComment` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `Settings`
  ADD PRIMARY KEY (`ID`),
  ADD UNIQUE KEY `SettingKey` (`SettingKey`);

ALTER TABLE `Settings`
  MODIFY `ID` int(11) NOT NULL AUTO_INCREMENT;
```

### Loans

```sql
CREATE TABLE `Loans` (
  `ID` int(11) NOT NULL,
  `Amount` float NOT NULL,
  `UserID` int(11) NOT NULL,
  `State` varchar(255) NOT NULL,
  `CreatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreatedBy` int(11) NOT NULL,
  `UpdatedAt` datetime DEFAULT NULL,
  `UpdatedBy` int(11) DEFAULT NULL,
  `Comment` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `Loans`
  ADD PRIMARY KEY (`ID`);

ALTER TABLE `Loans`
  MODIFY `ID` int(11) NOT NULL AUTO_INCREMENT;
```

### Transaction

```sql
CREATE TABLE `Transactions` (
  `ID` int(11) NOT NULL,
  `Amount` float NOT NULL,
  `Type` varchar(255) NOT NULL,
  `State` varchar(255) NOT NULL,
  `UserID` int(11) NOT NULL,
  `CreatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreatedBy` int(11) NOT NULL,
  `UpdatedAt` datetime DEFAULT NULL,
  `UpdatedBy` int(11) DEFAULT NULL,
  `UpdateComment` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `Transactions`
  ADD PRIMARY KEY (`ID`);

ALTER TABLE `Transactions`
  MODIFY `ID` int(11) NOT NULL AUTO_INCREMENT;
```

### Groups

```sql
CREATE TABLE `Groups` (
  `ID` int(11) NOT NULL,
  `GroupKey` varchar(255) NOT NULL,
  `Email` varchar(255) NOT NULL,
  `CreatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreatedBy` int(11) NOT NULL,
  `UpdatedAt` datetime DEFAULT NULL,
  `UpdatedBy` int(11) DEFAULT NULL,
  `UpdateComment` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `Groups`
  ADD PRIMARY KEY (`ID`),
  ADD UNIQUE KEY `GroupKey` (`GroupKey`);

ALTER TABLE `Groups`
  MODIFY `ID` int(11) NOT NULL AUTO_INCREMENT;
```

### Users

```sql
CREATE TABLE `Users` (
  `ID` int(11) NOT NULL,
  `State` varchar(255) NOT NULL,
  `FirstName` varchar(255) NOT NULL,
  `LastName` varchar(255) NOT NULL,
  `Email` varchar(255) NOT NULL,
  `Phone` varchar(255) NOT NULL,
  `Street` varchar(255) NOT NULL,
  `StreetNumber` varchar(255) NOT NULL,
  `Zip` varchar(255) NOT NULL,
  `Country` varchar(255) NOT NULL,
  `Birthday` date NOT NULL,
  `EntranceDate` date NOT NULL,
  `LeavingDate` date DEFAULT NULL,
  `AdditionalInfos` varchar(255) DEFAULT NULL,
  `LastActivityAt` datetime DEFAULT NULL,
  `CreatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `CreatedBy` int(11) NOT NULL,
  `UpdatedAt` datetime DEFAULT NULL,
  `UpdatedBy` int(11) DEFAULT NULL,
  `UpdateComment` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `Users`
  ADD PRIMARY KEY (`ID`),
  ADD UNIQUE KEY `Email` (`Email`);

ALTER TABLE `Users`
  MODIFY `ID` int(11) NOT NULL AUTO_INCREMENT;
```

### GroupUsers

```sql
CREATE TABLE `GroupUsers` (
  `ID` int(11) NOT NULL,
  `GroupID` int(11) NOT NULL,
  `UserID` int(11) NOT NULL,
  `IsLeader` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

ALTER TABLE `GroupUsers`
  ADD PRIMARY KEY (`ID`),
  ADD UNIQUE KEY `GroupID` (`GroupID`,`UserID`);

ALTER TABLE `GroupUsers`
  MODIFY `ID` int(11) NOT NULL AUTO_INCREMENT;
```

### api

Create user:

- POST to `/api/user/auth/create` with body:
  `{ "email": "tina@teewurst.org", "password": "**********", "lastname": "teewurst", "firstname": "tina", "mobile": "8348349", "street": "Fleischergasse", "credit_date": "2018-03-12" }`

Insert user transactions:

- POST to `/api/transaction` with body:
  `{ "transactions": [ { "amount": 100.00, "created_at": "2019-12-27 17:30", "category_id": 1, "status": 1 } ], "user": { "id": 176 } }`

Update doorcode:

- PATCH to `/api/settings/doorcode` with body:
  `{ "doorcode": "Außen: 225588 Innen:685259", "updated_at": "2019-12-23 14:00", "updated_by": 176 }`

Update user:

- PATCH to `/api/user` with body:
  `{ "id": 1, "username": "ro.ri", "email": "roland.rindfleisch@web.de", "lastname": "Rindfleisch", "firstname": "Roland", "mobile": "2837432847", "street": "Industriestraße 101", "zip": "04229", "city": "Leipzig", "date_of_birth": "1901-08-19", "date_of_entry": "2020-03-03" }`

for other routes see @ `routes.go`
