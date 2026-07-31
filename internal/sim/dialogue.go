package sim

func PlantDialogue(p *Plant) string {
	if !p.IsAlive() {
		return "..."
	}
	switch WaterStatus(p) {
	case "critical":
		return pick(p, []string{"I'm so thirsty...", "Please water me!", "Not doing well."})
	case "thirsty":
		return pick(p, []string{"Could use some water.", "Getting dry.", "A little thirsty."})
	case "too wet":
		return pick(p, []string{"A bit soggy...", "Too much water.", "Roots feel weird."})
	case "happy":
		return pick(p, []string{"Feeling good!", "All watered.", "Happy plant."})
	default:
		if p.Health < 30 {
			return pick(p, []string{"Not feeling great.", "Need some care.", "Wilting a little."})
		}
		return pick(p, []string{"Doing fine.", "Hanging in there.", "Ok for now."})
	}
}

func pick(p *Plant, lines []string) string {
	if len(lines) == 0 {
		return "..."
	}
	return lines[(int(p.Temperament)+len(p.Name))%len(lines)]
}
