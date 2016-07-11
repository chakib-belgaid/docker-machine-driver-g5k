package api

import (
    "encoding/json"
    "fmt"
    "time"
)

type Job struct {
    Uid       int       `json:"uid"`
    State     string    `json:"state"`
    Timelife  int       `json:"walltime"`
    Types     []string  `json:"types"`
    StartTime int       `json:"started_at"`
    Links     []Link    `json:"links"`
    Nodes     []string  `json:"assigned_nodes"`
}

func (a *Api) SubmitJob() (int, error) {
    var urlSubmit string = fmt.Sprintf("%s/sites/%s/jobs", G5kApiFrontend, a.Site)
    var job Job

    if resp, err := a.post(urlSubmit, `{"resources": "nodes=1,walltime=0:05:00", "command": "sleep 300", "types": ["deploy"]}`); err != nil {
        return 0, err
    } else {
        err = json.Unmarshal(resp, &job)
        return job.Uid, err
    }
}

// Refresh job's state
func (a *Api) GetJob(jobId int) (*Job, error) {
    job := new(Job)
    url := fmt.Sprintf("%s/sites/%s/jobs/%v", G5kApiFrontend, a.Site, jobId)

    if resp, err := a.get(url); err != nil {
        return &Job{}, err
    } else {
        err = json.Unmarshal(resp, &job)
        return job, err
    }
}

// Returns the job's current state
func (a *Api) GetJobState(jobId int) (string, error) {
    if job, err := a.GetJob(jobId); err != nil {
        return "", err
    } else if a.jobIsOver(job) {
        return "terminated", nil
    } else {
        return job.State, nil
    }
}

// Returns true if the job expired, false otherwise
func (a *Api) jobIsOver(job *Job) bool {
    currentTime := time.Now().Unix()
    startTime := int64(job.StartTime)
    timelife := int64(job.Timelife)

    if (currentTime - startTime) >= timelife {
        return true
    } else {
        return false
    }
}

// Free the nodes allocated to the jobs
func (a *Api) KillJob(jobId int) error {
    url := fmt.Sprintf("%s/sites/%s/jobs/%v", G5kApiFrontend, a.Site, jobId)

    _, err := a.del(url)

    return err
}

func (a *Api) waitJobIsReady(job *Job) bool {
    var err error
    tmp_job := new(Job)

    for job.State == "waiting" || job.State == "tolaunch" || job.State == "launching" {
        if tmp_job, err = a.GetJob(job.Uid); err != nil {
            return false
        }
        *job = *tmp_job
        time.Sleep(3*time.Second)
    }
    fmt.Println(job)

    // If the launching failed
    if job.State != "running" {
        return false
    } else {
        return true
    }
}
