import {ChildProcess, spawn} from 'child_process'
import gulp from 'gulp'
import path from 'path'

export function build(): ChildProcess {
    return spawn('yarn', ['-s', 'run', 'build'], {
        stdio: 'inherit',
        shell: true,
        env: {...process.env, NODE_OPTIONS: ''}
    })
}

export function watch(): ChildProcess {
    return spawn('yarn', ['-s', 'run', 'dev'], {
        stdio: 'inherit',
        shell: true,
        env: {...process.env, NODE_OPTIONS: '--max_old_space_size=8192'},
    })
}

const PHABRICATOR_EXTENSION_FILES = path.join(__dirname, './build/phabricator/dist/**')
const PHABRICATOR_ASSETS_DIRECTORY = path.join(__dirname, '../../ui/assets/extension')

export function phabricator(): NodeJS.ReadWriteStream {
    return gulp.src(PHABRICATOR_EXTENSION_FILES).pipe(gulp.dest(PHABRICATOR_ASSETS_DIRECTORY))
}
