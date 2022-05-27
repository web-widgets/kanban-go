package data

import (
	"time"
)

func dataDown(d *DAO) {
	d.mustExec("DELETE from cards")
	d.mustExec("DELETE from columns")
	d.mustExec("DELETE from rows")
	d.mustExec("DELETE from binary_data")
	d.mustExec("DELETE from users")
	d.mustExec("DELETE from assigned_users")
}

func formatDate(date string) *time.Time {
	t, _ := time.Parse(time.RFC3339, date)
	return &t
}

func dataUp(d *DAO) {
	db := d.GetDB()

	stage1 := Column{Name: "Backlog", Index: 1}
	db.Create(&stage1)
	stage2 := Column{Name: "In Progress", Index: 2}
	db.Create(&stage2)
	stage3 := Column{Name: "Testing", Index: 3}
	db.Create(&stage3)
	stage4 := Column{Name: "Done", Index: 4}
	db.Create(&stage4)

	row1 := Row{Name: "Feature", Index: 1}
	db.Create(&row1)
	row2 := Row{Name: "Task", Index: 2}
	db.Create(&row2)

	data1 := BinaryData{Name: "demo.png", Path: "x001"}
	db.Create(&data1)
	data2 := BinaryData{Name: "demo.png", Path: "x001"}
	db.Create(&data2)
	data3 := BinaryData{Name: "demo.png", Path: "x001"}
	db.Create(&data3)

	card1 := Card{
		Name:      "Integration with Angular/React",
		ColumnID:  stage1.ID,
		RowID:     row1.ID,
		Index:     1,
		Priority:  1,
		Color:     "#65D3B3",
		StartDate: formatDate("2018-01-01T00:00:00Z"),
	}
	db.Create(&card1)
	card2 := Card{
		Name:     "Archive the cards/boards",
		ColumnID: stage1.ID,
		RowID:    row1.ID,
		Index:    1,
		Priority: 3,
		Color:    "#58C3FE",
		Progress: 1,
	}
	db.Create(&card2)
	card3 := Card{
		Name:      "Searching and filtering",
		ColumnID:  stage1.ID,
		RowID:     row2.ID,
		Index:     2,
		Priority:  1,
		Color:     "#58C3FE",
		Progress:  1,
		StartDate: formatDate("2018-01-01T00:00:00Z"),
	}
	db.Create(&card3)
	card4 := Card{
		Name:         "Set the tasks priorities",
		ColumnID:     stage2.ID,
		RowID:        row1.ID,
		Color:        "#FFC975",
		Progress:     75,
		StartDate:    formatDate("2018-01-01T00:00:00Z"),
		AttachedData: []*BinaryData{&data1},
	}
	db.Create(&card4)
	card5 := Card{
		Name:      "Custom icons",
		ColumnID:  stage2.ID,
		RowID:     row2.ID,
		Color:     "#65D3B3",
		StartDate: formatDate("2019-01-01T00:00:00Z"),
	}
	db.Create(&card5)
	card6 := Card{
		Name:      "Integration with Gantt",
		ColumnID:  stage2.ID,
		RowID:     row2.ID,
		Color:     "#FFC975",
		Progress:  75,
		StartDate: formatDate("2020-01-01T00:00:00Z"),
	}
	db.Create(&card6)
	card7 := Card{
		Name:     "Drag and drop",
		ColumnID: stage3.ID,
		RowID:    row1.ID,
		Priority: 1,
		Color:    "#58C3FE",
		Progress: 100,
	}
	db.Create(&card7)
	card8 := Card{
		Name:         "Adding images",
		ColumnID:     stage3.ID,
		RowID:        row2.ID,
		Color:        "#58C3FE",
		AttachedData: []*BinaryData{&data1},
	}
	db.Create(&card8)
	card9 := Card{
		Name:      "Create cards and lists from the UI and from code",
		ColumnID:  stage4.ID,
		RowID:     row1.ID,
		Priority:  3,
		Color:     "#65D3B3",
		StartDate: formatDate("2018-06-08T00:00:00Z"),
	}
	db.Create(&card9)
	card10 := Card{
		Name:     "Draw swimlanes",
		ColumnID: stage4.ID,
		RowID:    row1.ID,
		Color:    "#FFC975",
	}
	db.Create(&card10)
	card11 := Card{
		Name:      "Progress bar",
		ColumnID:  stage4.ID,
		RowID:     row2.ID,
		Priority:  1,
		Color:     "#FFC975",
		Progress:  100,
		StartDate: formatDate("2018-01-01T00:00:00Z"),
	}
	db.Create(&card11)

	db.Create(&User{
		ID:     1,
		Avatar: "../assets/user.jpg",
		Name:   "Steve Smith",
	})
	db.Create(&User{
		ID:     2,
		Avatar: "../assets/user-2.jpg",
		Name:   "Aaron Long",
	})
	db.Create(&User{
		ID:     3,
		Avatar: "../assets/user-3.jpg",
		Name:   "Angela Allen",
	})
	db.Create(&User{
		ID:     4,
		Avatar: "../assets/user-4.jpg",
		Name:   "Angela Long",
	})
	db.Create(&User{
		ID:   5,
		Name: "John Doe",
	})
}
