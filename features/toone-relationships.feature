Feature: It should be possible to mutate with relationships

    Background: We have Johny user
        Given I send query:
            """
            mutation {
            deleteAllUsers
            deleteAllCompanies
            createUser(input:{id:"johny",firstName:"John",lastName:"Doe"}) { id }
            createCompany(input:{id:"test",name:"test company"}) { id }
            }
            """