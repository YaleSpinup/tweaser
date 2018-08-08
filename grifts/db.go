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

	nfq, err := newQuestion("What is the next feature you would like to see implemented?", "single", nf.ID, true)
	if err != nil {
		return err
	}

	if _, err = newAnswer("container service", "choice", nfq.ID); err != nil {
		return err
	}
	if _, err = newAnswer("standalone databases", "choice", nfq.ID); err != nil {
		return err
	}
	if _, err = newAnswer("serverless computing", "choice", nfq.ID); err != nil {
		return err
	}

	lb, err := newCampaign("Favorite Feature", time.Now().Add(72*time.Hour), time.Now().Add(144*time.Hour), true)
	if err != nil {
		return err
	}

	lbq, err := newQuestion("What is your favorite current feature?", "single", lb.ID, true)
	if err != nil {
		return err
	}

	if _, err = newAnswer("servers for regulated data", "choice", lbq.ID); err != nil {
		return err
	}
	if _, err = newAnswer("tryit", "choice", lbq.ID); err != nil {
		return err
	}
	if _, err = newAnswer("windows servers", "choice", lbq.ID); err != nil {
		return err
	}
	if _, err = newAnswer("external requests", "choice", lbq.ID); err != nil {
		return err
	}

	fd, err := newCampaign("Favorite Developer", time.Now().Add(-72*time.Hour), time.Now(), true)
	if err != nil {
		return err
	}

	fdq, err := newQuestion("Who is your favorite developer?", "single", fd.ID, true)
	if err != nil {
		return err
	}

	if _, err = newAnswer("Camden", "choice", fdq.ID); err != nil {
		return err
	}
	if _, err = newAnswer("Tenyo", "choice", fdq.ID); err != nil {
		return err
	}
	if _, err = newAnswer("Galen", "choice", fdq.ID); err != nil {
		return err
	}
	if _, err = newAnswer("Andrew", "choice", fdq.ID); err != nil {
		return err
	}

	dc, err := newCampaign("Disabled Campaign", time.Now().Add(-72*time.Hour), time.Now(), false)
	if err != nil {
		return err
	}

	dcq, err := newQuestion("How do you feel about disabled campaigns?", "single", dc.ID, true)
	if err != nil {
		return err
	}

	if _, err = newAnswer("Good", "choice", dcq.ID); err != nil {
		return err
	}
	if _, err = newAnswer("Bad", "choice", dcq.ID); err != nil {
		return err
	}
	if _, err = newAnswer("Meh", "choice", dcq.ID); err != nil {
		return err
	}

	mc, err := newCampaign("Multi Question", time.Now(), time.Now().Add(36*time.Hour), true)
	if err != nil {
		return err
	}

	mcq1, err := newQuestion("How do you feel about too many questions?", "single", mc.ID, true)
	if err != nil {
		return err
	}
	if _, err = newAnswer("Good", "choice", mcq1.ID); err != nil {
		return err
	}
	if _, err = newAnswer("Bad", "choice", mcq1.ID); err != nil {
		return err
	}
	if _, err = newAnswer("Meh", "choice", mcq1.ID); err != nil {
		return err
	}

	mcq2, err := newQuestion("How do you feel about multiple choice questions?", "multi", mc.ID, true)
	if err != nil {
		return err
	}
	if _, err = newAnswer("Good", "choice", mcq2.ID); err != nil {
		return err
	}
	if _, err = newAnswer("Bad", "choice", mcq2.ID); err != nil {
		return err
	}
	if _, err = newAnswer("Other", "input", mcq2.ID); err != nil {
		return err
	}

	mcq3, err := newQuestion("How do you feel about free form questions?", "input", mc.ID, false)
	if err != nil {
		return err
	}
	if _, err = newAnswer("", "input", mcq3.ID); err != nil {
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
	_, err = tx.ValidateAndSave(&campaign)
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

func newQuestion(text, qType string, campaignID uuid.UUID, enabled bool) (*models.Question, error) {
	tx, err := pop.Connect("development")
	if err != nil {
		return nil, err
	}

	question := models.Question{Text: text, CampaignID: campaignID, Enabled: enabled, Type: qType}
	_, err = tx.ValidateAndSave(&question)
	if err != nil {
		return nil, err
	}

	return &question, nil
}

func newAnswer(text, aType string, questionID uuid.UUID) (*models.Answer, error) {
	tx, err := pop.Connect("development")
	if err != nil {
		return nil, err
	}

	answer := models.Answer{Text: text, QuestionID: questionID, Type: aType}
	_, err = tx.ValidateAndSave(&answer)
	if err != nil {
		return nil, err
	}

	return &answer, nil
}
