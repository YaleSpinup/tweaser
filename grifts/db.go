package grifts

import (
	"log"
	"time"

	"git.yale.edu/spinup/tweaser/models"
	"github.com/gobuffalo/pop"
	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		tx, err := pop.Connect("development")
		if err != nil {
			return err
		}

		err = seedCampaigns(c, tx)
		if err != nil {
			log.Println(err)
		}

		return err
	})

})

func seedCampaigns(c *grift.Context, tx *pop.Connection) error {
	nextFeature := models.Campaign{Name: "Determine Feature Priority", StartDate: time.Now(), EndDate: time.Now().Add(72 * time.Hour), Enabled: true}
	_, err := tx.ValidateAndSave(&nextFeature)
	if err != nil {
		return err
	}

	nfQuestions := models.Question{Text: "What is the next feature you would like to see implemented?", CampaignID: nextFeature.ID, Enabled: true}
	_, err = tx.ValidateAndSave(&nfQuestions)
	if err != nil {
		return err
	}

	likeBest := models.Campaign{Name: "Favorite Feature", StartDate: time.Now().Add(72 * time.Hour), EndDate: time.Now().Add(144 * time.Hour), Enabled: true}
	_, err = tx.ValidateAndSave(&likeBest)
	if err != nil {
		return err
	}

	lbQuestions := models.Question{Text: "What is your favorite current feature?", CampaignID: likeBest.ID, Enabled: true}
	_, err = tx.ValidateAndSave(&lbQuestions)
	if err != nil {
		return err
	}

	favoriteDev := models.Campaign{Name: "Favorite Developer", StartDate: time.Now().Add(-72 * time.Hour), EndDate: time.Now(), Enabled: true}
	_, err = tx.ValidateAndSave(&favoriteDev)
	if err != nil {
		return err
	}

	fdQuestions := models.Question{Text: "Who is your favorite developer?", CampaignID: favoriteDev.ID, Enabled: true}
	_, err = tx.ValidateAndSave(&fdQuestions)
	if err != nil {
		return err
	}

	disabledCampaign := models.Campaign{Name: "Disabled Campaign", StartDate: time.Now().Add(-72 * time.Hour), EndDate: time.Now(), Enabled: false}
	_, err = tx.ValidateAndSave(&disabledCampaign)
	if err != nil {
		return err
	}

	dcQuestions := models.Question{Text: "How do you feel about disabled campaigns?", CampaignID: disabledCampaign.ID, Enabled: true}
	_, err = tx.ValidateAndSave(&dcQuestions)
	if err != nil {
		return err
	}

	multiQuestion := models.Campaign{Name: "Multi Question", StartDate: time.Now(), EndDate: time.Now().Add(36 * time.Hour), Enabled: true}
	_, err = tx.ValidateAndSave(&multiQuestion)
	if err != nil {
		return err
	}

	mqQuestions1 := models.Question{Text: "How do you feel about too many questions?", CampaignID: multiQuestion.ID, Enabled: true}
	_, err = tx.ValidateAndSave(&mqQuestions1)
	if err != nil {
		return err
	}

	mqQuestions2 := models.Question{Text: "How do you feel about way too many questions?", CampaignID: multiQuestion.ID, Enabled: true}
	_, err = tx.ValidateAndSave(&mqQuestions2)
	if err != nil {
		return err
	}

	mqQuestions3 := models.Question{Text: "How do you feel about way way too many questions?", CampaignID: multiQuestion.ID, Enabled: false}
	_, err = tx.ValidateAndSave(&mqQuestions3)
	if err != nil {
		return err
	}

	log.Println("Seeded Campaigns")
	return err
}
