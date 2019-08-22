package templates

var DummyModel = `type User {
	email: String @column(gorm: "type:varchar(64) comment '用户邮箱地址';default:null;") @validator(required: "true", type: "email")
	firstName: String
	lastName: String
	tasks: [Task!]! @relationship(inverse:"assignee")
}
type Task {
	title: String
	completed: Boolean
	dueDate: Time
	assignee: User @relationship(inverse:"tasks")
}
`