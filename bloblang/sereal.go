package bloblang

import (
	"errors"
	"fmt"

	"github.com/Weborama/Sereal/Go/sereal"
	"github.com/redpanda-data/benthos/v4/public/bloblang"
)

var ErrInvalidSerealVersion = errors.New("invalid sereal version")

func init() {
	err := bloblang.RegisterMethodV2("parse_sereal",
		bloblang.NewPluginSpec().
			Category("Parsing").
			Description("Parses a [Sereal](https://github.com/Sereal/Sereal) message into a structured document."),
		func(args *bloblang.ParsedParams) (bloblang.Method, error) {
			return bloblang.BytesMethod(func(data []byte) (any, error) {
				var sObj any
				if err := sereal.Unmarshal(data, &sObj); err != nil {
					return nil, err
				}
				return sObj, nil
			}), nil
		})
	if err != nil {
		panic(err)
	}

	err = bloblang.RegisterMethodV2("format_sereal",
		bloblang.NewPluginSpec().
			Category("Parsing").
			Description("Formats data as a [Sereal](https://github.com/Sereal/Sereal) message in bytes format.").
			Param(bloblang.NewInt64Param("version").
				Description("sereal encoder version to use.").
				Optional().Default(3)).
			Param(bloblang.NewBoolParam("perl_compat").
				Description("try to mimic Perl's structure as much as possible.").
				Optional().Default(false)).
			Param(bloblang.NewBoolParam("struct_as_map").
				Description("convert struct as map to save some bytes.").
				Optional().Default(false)).
			Param(bloblang.NewAnyParam("compression").
				Description("compress the main payload of the document using snappy, zlib, or zstd.").
				Optional().Default(nil)),
		func(args *bloblang.ParsedParams) (bloblang.Method, error) {
			var encoder *sereal.Encoder
			verOpt, err := args.GetOptionalInt64("version")
			if err != nil {
				return nil, err
			}
			if verOpt != nil {
				switch *verOpt {
				case 1:
					encoder = sereal.NewEncoder()
				case 2:
					encoder = sereal.NewEncoderV2()
				case 3:
					encoder = sereal.NewEncoderV3()
				default:
					return nil, fmt.Errorf("%w: %d", ErrInvalidSerealVersion, *verOpt)
				}
			} else {
				encoder = sereal.NewEncoderV3()
			}

			compressionOpt, err := args.Get("compression") // interface to allow later object
			if err != nil {
				return nil, err
			}
			if compressionOpt != nil {
				switch val := compressionOpt.(type) {
				case string: // use default value
					switch val {
					case "snappy":
						encoder.Compression = sereal.SnappyCompressor{Incremental: true}
					case "zlib":
						encoder.Compression = sereal.ZlibCompressor{}
					case "zstd": // need cgo
						encoder.Compression = sereal.ZstdCompressor{}
					}
					// TODO add struct matching for more compression option
				}
			}

			perlOpt, err := args.GetOptionalBool("perl_compat")
			if err != nil {
				return nil, err
			}
			if perlOpt != nil {
				encoder.PerlCompat = *perlOpt
			}

			structOpt, err := args.GetOptionalBool("struct_as_map")
			if err != nil {
				return nil, err
			}
			if structOpt != nil {
				encoder.StructAsMap = *structOpt
			}

			return func(v any) (any, error) {
				return encoder.Marshal(v)
			}, nil
		})
	if err != nil {
		panic(err)
	}
}
