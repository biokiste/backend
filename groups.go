package main

// GetGroupsWithUsers return groups with user ids and ids of group leaders
func (h Handlers) GetGroupsWithUsers() ([]Group, error) {
	results, err := h.DB.Query(`
		SELECT
			group_id, user_id, position_id
		FROM
			groups_users
		WHERE
			active=1
	`)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var groups []Group

	for results.Next() {
		var entry GroupUserEntry
		err = results.Scan(&entry.GroupID, &entry.UserID, &entry.PositionID)

		if err != nil {
			panic(err.Error())
		}

		var idx int = -1

		for i, g := range groups {
			if g.ID == entry.GroupID {
				idx = i
				break
			}
		}

		if idx == -1 {
			var newGroup Group
			newGroup.ID = entry.GroupID
			newGroup.UserIDs = append(newGroup.UserIDs, entry.UserID)
			if entry.PositionID == 1 {
				newGroup.LeaderIDs = append(newGroup.LeaderIDs, entry.UserID)
			}
			groups = append(groups, newGroup)
		} else {
			var group = groups[idx]
			group.UserIDs = append(group.UserIDs, entry.UserID)
			if entry.PositionID == 1 {
				group.LeaderIDs = append(group.LeaderIDs, entry.UserID)
			}
			groups[idx] = group
		}

	}

	return groups, err
}
