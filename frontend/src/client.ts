import axios from "axios"

class Client {

    client = axios.create({
        baseURL: "./api/",
        maxRedirects: 0
    })

    async getJob(name: string): Promise<any> {
        let job = await this.client.get(`jobs/${name}`)
        return job.data;
    }

    async checkJob(name: string) {
        await this.client.post(`jobs/${name}/check`)
    }

    async testJobActions(name: string) {
        await this.client.post(`jobs/${name}/test-actions`)
    }

    async listJobs(): Promise<any[]> {
        let jobs = await this.client.get("jobs")
        return jobs.data;
    }
}

export default new Client();