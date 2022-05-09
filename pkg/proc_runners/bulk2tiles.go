package proc_runners

import (
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog/log"
	"glacierpeak.app/openland/pkg/utils"
)

func Bulk2Tiles(dir string, outDir string, workers int, zoom int) {
	sources, _ := utils.WalkMatch(dir, "*.tif")
	log.Warn().Msgf("Running with %v workers", workers)

	for idx, source := range sources {
		sourceQuoted := "'" + source + "'"
		basePath := filepath.Dir(source)
		fileWithoutExt := utils.StripExt(source)
		vrtFile := filepath.Join(basePath, fileWithoutExt+".vrt")
		vrtFileQuoted := "'" + vrtFile + "'"
		// gdal_translate -of vrt -expand rgba /mnt/mapvault/Benchmark/TIF/Arizona Landscape Maps IC300 GeoTiff.tif temp.vrt
		out1, err := utils.RunCommand2(true, false, "gdal_translate", "-of", "vrt", "-expand", "rgba", sourceQuoted, vrtFileQuoted)
		log.Err(err).Msg(out1)
		out2, err := utils.RunCommand2(true, false, "gdal2tiles.py", "--zoom="+strconv.Itoa(zoom), "--processes="+strconv.Itoa(workers), "--xyz", "--resume", vrtFileQuoted, filepath.Join(outDir, strconv.Itoa(idx)))
		log.Err(err).Msg(out2)
	}
	log.Info().Msg("Done!")
}
