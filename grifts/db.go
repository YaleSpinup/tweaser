package grifts

import (
	"log"
	"time"

	"git.yale.edu/spinup/tweaser/models"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		err := seedCampaigns(c)
		if err != nil {
			log.Println(err)
		}

		return err
	})

})

func seedCampaigns(c *grift.Context) error {
	nf, err := newCampaign("Determine Feature Priority", time.Now(), time.Now().Add(72*time.Hour), true)
	if err != nil {
		return err
	}

	nfq, err := newQuestion("What is the next feature you would like to see implemented?", nf.ID, true)
	if err != nil {
		return err
	}

	if _, err = newAnswer("container service", nfq.ID); err != nil {
		return err
	}
	if _, err = newAnswer("standalone databases", nfq.ID); err != nil {
		return err
	}
	if _, err = newAnswer("serverless computing", nfq.ID); err != nil {
		return err
	}

	lb, err := newCampaign("Favorite Feature", time.Now().Add(72*time.Hour), time.Now().Add(144*time.Hour), true)
	if err != nil {
		return err
	}

	lbq, err := newQuestion("What is your favorite current feature?", lb.ID, true)
	if err != nil {
		return err
	}

	if _, err = newAnswer("servers for regulated data", lbq.ID); err != nil {
		return err
	}
	if _, err = newAnswer("tryit", lbq.ID); err != nil {
		return err
	}
	if _, err = newAnswer("windows servers", lbq.ID); err != nil {
		return err
	}
	if _, err = newAnswer("external requests", lbq.ID); err != nil {
		return err
	}

	fd, err := newCampaign("Favorite Developer", time.Now().Add(-72*time.Hour), time.Now(), true)
	if err != nil {
		return err
	}

	_, err = newQuestion("Who is your favorite developer?", fd.ID, true)
	if err != nil {
		return err
	}

	dc, err := newCampaign("Disabled Campaign", time.Now().Add(-72*time.Hour), time.Now(), false)
	if err != nil {
		return err
	}

	_, err = newQuestion("How do you feel about disabled campaigns?", dc.ID, true)
	if err != nil {
		return err
	}

	mq, err := newCampaign("Multi Question", time.Now(), time.Now().Add(36*time.Hour), true)
	if err != nil {
		return err
	}

	_, err = newQuestion("How do you feel about too many questions?", mq.ID, true)
	if err != nil {
		return err
	}

	_, err = newQuestion("How do you feel about way too many questions?", mq.ID, true)
	if err != nil {
		return err
	}

	_, err = newQuestion("How do you feel about way way too many questions?", mq.ID, false)
	if err != nil {
		return err
	}

	log.Println("Seeded Campaigns")
	return err
}

func newCampaign(name string, start, end time.Time, enabled bool) (*models.Campaign, error) {
	tx, err := pop.Connect("development")
	if err != nil {
		return nil, err
	}

	campaign := models.Campaign{Name: name, StartDate: start, EndDate: end, Enabled: enabled}
	out, err := tx.ValidateAndSave(&campaign)
	if err != nil {
		return nil, err
	}
	log.Println(out)
	return &campaign, nil
}

func newQuestion(text string, campaignID uuid.UUID, enabled bool) (*models.Question, error) {
	tx, err := pop.Connect("development")
	if err != nil {
		return nil, err
	}

	question := models.Question{Text: text, CampaignID: campaignID, Enabled: enabled}
	_, err = tx.ValidateAndSave(&question)
	if err != nil {
		return nil, err
	}

	return &question, nil
}

func newAnswer(text string, questionID uuid.UUID) (*models.Answer, error) {
	tx, err := pop.Connect("development")
	if err != nil {
		return nil, err
	}

	answer := models.Answer{Text: text, QuestionID: questionID}
	_, err = tx.ValidateAndSave(&answer)
	if err != nil {
		return nil, err
	}

	return &answer, nil
}
