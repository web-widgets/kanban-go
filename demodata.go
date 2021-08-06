package main

func dataDown() {
	mustExec("DELETE from cards")
	mustExec("DELETE from stages")
	mustExec("DELETE from binary_data")
}

func dataUp() {
	stage1 := Stage{Name: "ToDo"}
	db.Create(&stage1)
	stage2 := Stage{Name: "In Progress"}
	db.Create(&stage2)
	stage3 := Stage{Name: "Testing"}
	db.Create(&stage3)
	stage4 := Stage{Name: "Done"}
	db.Create(&stage4)

	data1 := BinaryData{Name: "demo.png", Path: "x001"}
	db.Create(&data1)
	data2 := BinaryData{Name: "demo.png", Path: "x001"}
	db.Create(&data2)
	data3 := BinaryData{Name: "demo.png", Path: "x001"}
	db.Create(&data3)

	card1 := Card{
		Name:         "Reordering in the Kanban",
		StageID:      stage4.ID,
		AttachedData: []*BinaryData{&data1},
		Index:        1,
	}
	db.Create(&card1)
	card2 := Card{
		Name:    "UX optimization",
		StageID: stage2.ID,
		Index:   1,
	}
	db.Create(&card2)
	card3 := Card{
		Name:    "Accessibility",
		StageID: stage2.ID,
		Index:   2,
	}
	db.Create(&card3)
}

func mustExec(sql string) {
	err := db.Exec(sql).Error
	if err != nil {
		panic(err)
	}
}
