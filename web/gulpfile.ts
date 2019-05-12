import log from 'fancy-log'
import {Stats} from 'webpack'


const WEBPACK_STATS_OPTIONS: Stats.ToStringOptions & { colors?: boolean } = {}

const logWebpackStats = (stats: Stats) => log(stats.toString(WEBPACK_STATS_OPTIONS))

export async function webpack(): Promise<void> {

}
