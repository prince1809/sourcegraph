import {ChildProcess, spawn} from 'child_process';

export function build(): ChildProcess {
    return spawn('yarn', ['-s', 'run','build'], {
        stdio: 'inherit',
        shell: true,
        env: { ...process.env, NODE_OPTIONS: ''}
    })
}
