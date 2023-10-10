package models_test

import (
	"testing"
)

func TestMatch_AddTeam(t *testing.T) {

	//TODO we don't have this on match anymore
	//REVAMP

	// testutils.InitializeTestStructs();
	// _, _, db := testutils.GetTestingStructs()

	// tx, err := db.Begin()
	// if err != nil {
	// 	t.Fatalf("Could not create test transaction, %s", err)
	// }
	// defer tx.Commit()

	// matchID := "008330f8-2e79-4161-8924-4be4e0006bdf"
	// match := models.Match{ID: repository.NewSQLCompatUUIDFromStr(matchID)}
	// team := models.Team{MatchID: repository.NewSQLCompatUUIDFromStr(matchID)}
	// err = match.AddTeam(team, db)
	// if(err != nil){
	// 	t.Errorf("Unexpected error %s", err.Error())
	// }

	// err = db.Select(&match)
	// if err != nil {
	// 	t.Errorf("Error selecting match %s", err.Error())
	// }

	// t.Logf("Record Found %+v", &match)

}
