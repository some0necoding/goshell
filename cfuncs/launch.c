#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>
#include <launch.h>

int lsh_launch(char **args) {

    pid_t pid, wpid;
    int status;

    pid = fork(); 
    
    // on successfull fork 0 is returned to the newly created child process, 
    // while child's pid is returned to the parent process. Also note that when
    // fork is called both processes (child and parent) are running the same program,
    // so the child (pid = 0) will exec (overriding the program running at that time), 
    // while the parent (pid = child's pid) will wait for the child process to terminate 
    
    if ((pid == 0) && (execvp(args[0], args) == -1)) {
        return -1;
    } else if (pid < 0) {
        return -1;
    } else {
        do {
            wpid = waitpid(pid, &status, WUNTRACED);
        } while (!WIFEXITED(status) && !WIFSIGNALED(status));   
    }

    return 1;
}
