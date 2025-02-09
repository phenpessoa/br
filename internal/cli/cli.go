package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/phenpessoa/br"
)

func Run() error {
	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprint(flags.Output(), usage)
		flags.PrintDefaults()
	}

	var flagSilent bool
	flags.BoolVar(
		&flagSilent,
		"silent",
		false,
		"If set to true, no output will be produced. Instead, br will exit with code 0 if the provided document is valid or code 1 if it is invalid.\n"+
			"This flag is ignored if generating a valid document.\nIf br is used incorrectly, it will return exit code 2.",
	)

	// var numericCNPJ bool
	// flags.BoolVar(
	// 	&numericCNPJ,
	// 	"ncpj",
	// 	false,
	// 	"When validating a CNPJ, if set to true, only numeric CNPJs will be considered valid and alpha numeric CNPJs will be considered invalid, even if they are actually valid.\n"+
	// 		"When generating a CNPJ, if set to true, only nymeric CNPJs will be generated.",
	// )

	if err := flags.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	args := flags.Args()
	switch len(args) {
	case 1:
		cmd := args[0]
		switch strings.ToLower(cmd) {
		case "cpf":
			cpf := br.GenerateCPF()
			fmt.Println(cpf)
			return nil
		case "cnpj":
			cnpj := br.GenerateCNPJ()
			fmt.Println(cnpj)
			return nil
		default:
			if flagSilent {
				os.Exit(2)
			}
			fmt.Printf("unknown document passed: %s\n", cmd)
			flags.Usage()
			return nil
		}
	case 2:
		cmd := strings.TrimSpace(args[0])
		arg := strings.TrimSpace(args[1])
		switch strings.ToLower(cmd) {
		case "cpf":
			cpf, err := br.NewCPF(arg)
			if err != nil {
				if flagSilent {
					os.Exit(1)
				}
				fmt.Printf("CPF %s is invalid 🚫\n", arg)
				return nil
			}
			if flagSilent {
				os.Exit(0)
			}
			fmt.Printf("CPF %s is valid ✅\n", cpf)
			return nil
		case "cnpj":
			cnpj, err := br.NewCNPJ(arg)
			if err != nil {
				if flagSilent {
					os.Exit(1)
				}
				fmt.Printf("CNPJ %s is invalid 🚫\n", arg)
				return nil
			}
			if flagSilent {
				os.Exit(0)
			}
			fmt.Printf("CNPJ %s is valid ✅\n", cnpj)
			return nil
		default:
			if flagSilent {
				os.Exit(2)
			}
			fmt.Printf("unknown document passed: %s\n", cmd)
			flags.Usage()
			return nil
		}
	default:
		if flagSilent {
			os.Exit(2)
		}
		flags.Usage()
		return nil
	}
}
