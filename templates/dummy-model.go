package templates

var DummyModel = `type User @entity {
	email: String @column(gorm: "type:varchar(64) comment '用户邮箱地址';default:null;") @validator(required: "true", type: "email")
	firstName: String @column
	lastName: String @column

	tasks: [Task!]! @relationship(inverse:"assignee")
}

type Task @entity {
	title: String @column
	completed: Boolean @column
	dueDate: Time @column

	assignee: User @relationship(inverse:"tasks")
}

`
