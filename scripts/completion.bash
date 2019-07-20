#/usr/bin/env bash

_dynconf_completion() {
	words=${#COMP_WORDS[@]}
	if [ $words -le 2 ]; then
		COMPREPLY=($(compgen -W "apply check show help version" -- "${COMP_WORDS[1]}"))
	elif [ $words -eq 3 ]; then
		subcommand="${COMP_WORDS[1]}"
		if [ "$subcommand" == "apply" ] || [ "$subcommand" == "check" ] || [ "$subcommand" == "show" ]; then
			compopt -o filenames
			COMPREPLY=($(compgen -f -X "!*.yml" -- "${COMP_WORDS[2]}") $(compgen -d -- "${COMP_WORDS[2]}"))
		fi
	fi
}

complete -F _dynconf_completion dynconf
