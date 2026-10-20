package menus

import (
	"strings"
)

func CollateMenuOptions(uniquekey string, menuoptions []string, cli *map[string]Menu) {
	combined := GetMenuMap(uniquekey, menuoptions)
	menu := Menu{menuname: uniquekey}
	menuopts := make([]MenuOption, len(combined[uniquekey]))

	for i := 0; i < len(menuoptions); i++ {
		splitOpt := strings.Split(menuoptions[i], " ")
		menuopts[i] = MenuOption{
			accessor:      splitOpt[0],
			value:         splitOpt[1],
			uniquemenukey: splitOpt[2],
			action:        splitOpt[3],
		}
	}
	menu.menuoptions = menuopts
	(*cli)[uniquekey] = menu
}

func collate() (map[string]Menu, map[string]func(*Runtime)) {
	allmenukeys := []string{
		"main",
		"main.state",
		"main.pipeline",
		"main.sources",
		"main.destinations",
		"main.transforms",
		"main.sequences",
	}

	cli := GetCliMap(allmenukeys)

	mainopts := []string{
		"q quit main.quit exit",
		"1 state main.state menu",
		"2 pipelines main.pipeline menu",
		"3 sources main.sources menu",
		"4 destinations main.destinations menu",
		"5 transforms main.transforms menu",
		"6 sequences main.sequences menu",
	}

	stateops := []string{
		"q back main.state.quit exit",
		"1 show main.state.show action",
		"2 set main.state.set action",
	}

	pipelineops := []string{
		"q back main.pipeline.quit exit",
		"1 new main.pipeline.newpipeline action",
		"2 show main.pipeline.showpipeline action",
		"3 delete main.pipeline.deletepipeline action",
	}

	sourcesops := []string{
		"q back main.sources.quit exit",
		"1 new main.sources.newsource action",
		"2 show main.sources.showsource action",
		"3 read main.sources.readsource action",
		"4 delete main.sources.deletesource action",
	}

	destops := []string{
		"q back main.destinations.quit exit",
	}

	transformsops := []string{
		"q back main.transforms.quit exit",
	}

	sequencesops := []string{
		"q back main.sequences.quit exit",
	}

	/*
		IMPORTANT, when adding a new menu's options they need to be in the
		correct order with the allmenukeys variable above, otherwise
		the menus will be mismatched
	*/

	allmenuoptions := [][]string{
		mainopts,
		stateops,
		pipelineops,
		sourcesops,
		destops,
		transformsops,
		sequencesops,
	}

	for i := 0; i < len(allmenuoptions); i++ {
		CollateMenuOptions(allmenukeys[i], allmenuoptions[i], &cli)
	}

	actions := make(map[string]func(*Runtime))
	actions["main.pipeline.newpipeline"] = WriteToPipelines
	actions["main.pipeline.showpipeline"] = ShowPipelines
	actions["main.pipeline.deletepipeline"] = DeletePipeline
	actions["main.sources.newsource"] = NewSource
	actions["main.sources.showsource"] = ShowSources
	actions["main.sources.readsource"] = ReadSource
	actions["main.sources.deletesource"] = DeleteSource
	actions["main.state.show"] = ShowState
	actions["main.state.set"] = SetState

	return cli, actions
}
