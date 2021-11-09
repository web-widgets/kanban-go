package data

func dataDown() {
	mustExec("DELETE from cards")
	mustExec("DELETE from columns")
	mustExec("DELETE from rows")
	mustExec("DELETE from stages")
	mustExec("DELETE from binary_data")
}

func dataUp() {
	stage1 := Column{Name: "ToDo"}
	db.Create(&stage1)
	stage2 := Column{Name: "In Progress"}
	db.Create(&stage2)
	stage3 := Column{Name: "Testing"}
	db.Create(&stage3)
	stage4 := Column{Name: "Done"}
	db.Create(&stage4)

	row1 := Row{Name: "Feature"}
	db.Create(&row1)
	row2 := Row{Name: "Task"}
	db.Create(&row2)

	data1 := BinaryData{Name: "demo.png", Path: "x001"}
	db.Create(&data1)
	data2 := BinaryData{Name: "demo.png", Path: "x001"}
	db.Create(&data2)
	data3 := BinaryData{Name: "demo.png", Path: "x001"}
	db.Create(&data3)

	card1 := Card{
		Name:         "Reordering in the Kanban",
		ColumnID:     stage4.ID,
		RowID:        row1.ID,
		AttachedData: []*BinaryData{&data1},
		Index:        1,
	}
	db.Create(&card1)
	card2 := Card{
		Name:     "UX optimization",
		ColumnID: stage2.ID,
		RowID:    row1.ID,
		Index:    1,
	}
	db.Create(&card2)
	card3 := Card{
		Name:     "Accessibility",
		ColumnID: stage2.ID,
		RowID:    row1.ID,
		Index:    2,
	}
	db.Create(&card3)
}
