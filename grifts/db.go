package grifts

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/YaleSpinup/tweaser/models"
	"github.com/gobuffalo/grift/grift"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
)

var _ = grift.Namespace("db", func() {
	_ = grift.Desc("seed", "Seeds a database")
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

	answerIDs := []uuid.UUID{}
	a, err := newAnswer("container service", nfq.ID, true)
	if err != nil {
		return err
	}
	answerIDs = append(answerIDs, a.ID)

	a, err = newAnswer("standalone databases", nfq.ID, true)
	if err != nil {
		return err
	}
	answerIDs = append(answerIDs, a.ID)

	a, err = newAnswer("serverless computing", nfq.ID, true)
	if err != nil {
		return err
	}
	answerIDs = append(answerIDs, a.ID)

	if len(c.Args) > 0 && c.Args[0] == "all" {
		times := rand.Intn(100)
		for i := 0; i <= times; i += 1 {
			userid := fmt.Sprintf("someuser%d", rand.Intn(100))
			answerID := answerIDs[rand.Intn(len(answerIDs))]
			_, err = newResponse(userid, "", nfq.ID, []uuid.UUID{answerID})
			if err != nil {
				return err
			}
		}
	}

	lb, err := newCampaign("Favorite Feature", time.Now().Add(72*time.Hour), time.Now().Add(144*time.Hour), true)
	if err != nil {
		return err
	}

	lbq, err := newQuestion("What is your favorite current feature?", "single", lb.ID, true)
	if err != nil {
		return err
	}

	answerIDs = []uuid.UUID{}
	if a, err = newAnswer("servers for regulated data", lbq.ID, true); err != nil {
		return err
	}
	answerIDs = append(answerIDs, a.ID)

	if a, err = newAnswer("tryit", lbq.ID, true); err != nil {
		return err
	}
	answerIDs = append(answerIDs, a.ID)

	if a, err = newAnswer("windows servers", lbq.ID, true); err != nil {
		return err
	}
	answerIDs = append(answerIDs, a.ID)

	if a, err = newAnswer("external requests", lbq.ID, true); err != nil {
		return err
	}
	answerIDs = append(answerIDs, a.ID)

	if len(c.Args) > 0 && c.Args[0] == "all" {
		times := rand.Intn(100)
		for i := 0; i <= times; i += 1 {
			userid := fmt.Sprintf("someuser%d", rand.Intn(100))
			answerID := answerIDs[rand.Intn(len(answerIDs))]
			_, err = newResponse(userid, "", lbq.ID, []uuid.UUID{answerID})
			if err != nil {
				return err
			}
		}
	}

	fd, err := newCampaign("Favorite Developer", time.Now().Add(-72*time.Hour), time.Now(), true)
	if err != nil {
		return err
	}

	fdq, err := newQuestion("Who is your favorite developer?", "single", fd.ID, true)
	if err != nil {
		return err
	}

	answerIDs = []uuid.UUID{}
	if a, err = newAnswer("Camden", fdq.ID, true); err != nil {
		return err
	}
	answerIDs = append(answerIDs, a.ID)

	if a, err = newAnswer("Tenyo", fdq.ID, true); err != nil {
		return err
	}
	answerIDs = append(answerIDs, a.ID)

	if a, err = newAnswer("Galen", fdq.ID, true); err != nil {
		return err
	}
	answerIDs = append(answerIDs, a.ID)

	if a, err = newAnswer("Andrew", fdq.ID, true); err != nil {
		return err
	}
	answerIDs = append(answerIDs, a.ID)

	if len(c.Args) > 0 && c.Args[0] == "all" {
		times := rand.Intn(100)
		for i := 0; i <= times; i += 1 {
			userid := fmt.Sprintf("someuser%d", rand.Intn(100))
			answerID := answerIDs[rand.Intn(len(answerIDs))]
			_, err = newResponse(userid, "", lbq.ID, []uuid.UUID{answerID})
			if err != nil {
				return err
			}
		}
	}

	dc, err := newCampaign("Disabled Campaign", time.Now().Add(-72*time.Hour), time.Now(), false)
	if err != nil {
		return err
	}

	dcq, err := newQuestion("How do you feel about disabled campaigns?", "single", dc.ID, true)
	if err != nil {
		return err
	}

	if _, err = newAnswer("Good", dcq.ID, true); err != nil {
		return err
	}
	if _, err = newAnswer("Bad", dcq.ID, true); err != nil {
		return err
	}
	if _, err = newAnswer("Meh", dcq.ID, true); err != nil {
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
	if _, err = newAnswer("Good", mcq1.ID, true); err != nil {
		return err
	}
	if _, err = newAnswer("Bad", mcq1.ID, true); err != nil {
		return err
	}
	if _, err = newAnswer("Meh", mcq1.ID, true); err != nil {
		return err
	}

	mcq2, err := newQuestion("How do you feel about multiple choice questions?", "multi", mc.ID, true)
	if err != nil {
		return err
	}
	if _, err = newAnswer("Good", mcq2.ID, true); err != nil {
		return err
	}
	if _, err = newAnswer("Bad", mcq2.ID, true); err != nil {
		return err
	}
	if _, err = newAnswer("Other", mcq2.ID, true); err != nil {
		return err
	}
	if _, err = newAnswer("Supercalifragilisticexpialidocious", mcq2.ID, false); err != nil {
		return err
	}

	_, err = newQuestion("How do you feel about disabled free form questions?", "input", mc.ID, false)
	if err != nil {
		return err
	}

	_, err = newQuestion("How do you feel about enabled free form questions?", "input", mc.ID, true)
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

func newAnswer(text string, questionID uuid.UUID, enabled bool) (*models.Answer, error) {
	tx, err := pop.Connect("development")
	if err != nil {
		return nil, err
	}

	answer := models.Answer{Text: text, QuestionID: questionID, Enabled: enabled}
	_, err = tx.ValidateAndSave(&answer)
	if err != nil {
		return nil, err
	}

	return &answer, nil
}

func newResponse(user, text string, questionID uuid.UUID, answerIDs []uuid.UUID) (*models.Response, error) {
	tx, err := pop.Connect("development")
	if err != nil {
		return nil, err
	}

	var answers models.Answers
	for _, a := range answerIDs {
		ans := models.Answer{}
		if err := tx.Find(&ans, a); err != nil {
			return nil, err
		}
		answers = append(answers, ans)
	}

	// create response
	response := models.Response{
		UserID:     user,
		Text:       text,
		QuestionID: questionID,
		Answers:    answers,
	}

	_, err = tx.ValidateAndSave(&response)
	if err != nil {
		return nil, err
	}

	// create response/answer association for each
	for _, a := range answerIDs {
		log.Println("Creating associations for response id", response, "to answer id", a)
		responseAnswer := models.ResponseAnswer{
			ResponseID: response.ID,
			AnswerID:   a,
			QuestionID: questionID,
		}
		log.Println("responseanswer", responseAnswer)

		_, err := tx.ValidateAndSave(&responseAnswer)
		if err != nil {
			return nil, err
		}
	}

	return &response, nil
}
