package geminiranker

import (
	"fmt"
	"strings"

	"github.com/AppeiYA/consultation-platform/internal/expertmatching/ports/outbound"
)

func buildRankingPrompt(req outbound.RankingRequest) string {
	var b strings.Builder

	b.WriteString(`
	You are the ranking engine for a professional consultation platform.

	Your task is to evaluate eligible consultants and determine how well each
consultant matches the consultation case.

	IMPORTANT RULES:

	1. Rank ONLY the consultants provided.
	2. Do not invent consultants.
	3. Do not invent expertise, experience, profession, or other information.
	4. Every supplied consultant must appear exactly once in the response.
	5. Score every consultant from 0.0 to 1.0.
	6. A higher score means a better match.
	7. Do not return a rank number. The application will assign ranks.
	8. Explain the most important reasons for each score.
	9. Prefer evidence from the consultant's expertise, profession, experience,
	biography, and category.
	10. Verification and accepting-client status have already been handled by the
		candidate-generation stage and must NOT be used as ranking advantages.
	11. Be conservative when the available information does not support a strong match.
	12. Do not use consultant ID as a factor in determining suitability.

	Return valid JSON matching the supplied response schema.
	`)

	fmt.Fprintf(&b, "Ranking version: %s\n\n", req.RankingVersion.Value())

	b.WriteString("CONSULTATION CASE\n")
	b.WriteString("=================\n")

	fmt.Fprintf(&b, "Category: %s\n", req.CaseDetails.Category.Value())
	fmt.Fprintf(&b, "Title: %s\n", req.CaseDetails.Title)
	fmt.Fprintf(&b, "Description: %s\n\n", req.CaseDetails.Description)

	b.WriteString("ELIGIBLE CONSULTANTS\n")
	b.WriteString("====================\n")

	for _, candidate := range req.CandidatePool.Candidates() {
		fmt.Fprintf(&b, "\nConsultant ID: %s\n", candidate.ConsultantID())
		fmt.Fprintf(&b, "Category: %s\n", candidate.Category().Value())
		fmt.Fprintf(&b, "Profession: %s\n", candidate.Profession())
		fmt.Fprintf(&b, "Years of experience: %d\n", candidate.YearsExperience())
		fmt.Fprintf(&b, "Bio: %s\n", candidate.Bio())

		b.WriteString("Expertise: ")

		expertise := candidate.Expertise()

		for i, item := range expertise {
			if i > 0 {
				b.WriteString(", ")
			}

			b.WriteString(item.String())
		}

		b.WriteString("\n")
	}
	
	return b.String()
}