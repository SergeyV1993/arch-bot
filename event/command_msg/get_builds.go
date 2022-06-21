package command_msg

import (
	"ArchitectureBot/repository"
	"fmt"
)

func CreateBuildMsg(build repository.ArchitectBuilds) string {
	text := fmt.Sprintf("*%s*", build.Name) + "\n" +
		fmt.Sprintf("Приблизительное расстояние к достопримечательности: _%.1f_ м.", build.Distance) + "\n" +
		fmt.Sprintf("[%s](%s)", build.Address, build.LinkMapAddress) + "\n" +
		fmt.Sprintf("*Описание: *_%s_", build.Description) + "\n" +
		fmt.Sprintf("[Более подробно по ссылке](%s)", build.Link) +
		"\n" + "\n"

	return text
}
