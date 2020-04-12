import axios from "axios"

class Client {

    client = axios.create({
        baseURL: "./api/",
        maxRedirects: 0
    })

    async getJob(id: string): Promise<any> {
        let job = await this.client.get(`jobs/${id}`)
        return job.data;
    }

    async checkJob(id: string) {
        await this.client.post(`jobs/${id}/check`)
    }

    async testJobActions(id: string) {
        await this.client.post(`jobs/${id}/test-actions`)
    }

    async listJobs(): Promise<any[]> {
        let jobs = await this.client.get("jobs")
        return jobs.data;
    }
}

export default new Client();