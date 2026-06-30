import request from './request'

export const listJobs = (params) => {
  return request.get('/jobs', { params })
}

export const getJob = (id) => {
  return request.get(`/jobs/${id}`)
}

export const createJob = (data) => {
  return request.post('/jobs', data)
}

export const updateJob = (id, data) => {
  return request.put(`/jobs/${id}`, data)
}

export const listApplications = (params) => {
  return request.get('/applications', { params })
}

export const getApplication = (id) => {
  return request.get(`/applications/${id}`)
}

export const createApplication = (data) => {
  return request.post('/applications', data)
}

export const updateApplicationStatus = (id, status) => {
  return request.put(`/applications/${id}/status`, { status })
}

export const listMessages = (applicationId) => {
  return request.get(`/messages/application/${applicationId}`)
}

export const sendMessage = (data) => {
  return request.post('/messages', data)
}

export const getStats = () => {
  return request.get('/stats')
}
