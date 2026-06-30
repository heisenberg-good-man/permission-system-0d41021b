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

export const listCandidates = (params) => {
  return request.get('/candidates', { params })
}

export const getCandidate = (params) => {
  return request.get('/candidates/detail', { params })
}

export const updateCandidateStatus = (data) => {
  return request.put('/candidates/status', data)
}

export const listNotes = (params) => {
  return request.get('/notes', { params })
}

export const createNote = (data) => {
  return request.post('/notes', data)
}

export const listInterviews = (params) => {
  return request.get('/interviews', { params })
}

export const getInterview = (id) => {
  return request.get(`/interviews/${id}`)
}

export const createInterview = (data) => {
  return request.post('/interviews', data)
}

export const rescheduleInterview = (id, data) => {
  return request.put(`/interviews/${id}/reschedule`, data)
}

export const cancelInterview = (id) => {
  return request.put(`/interviews/${id}/cancel`)
}

export const completeInterview = (id, data) => {
  return request.put(`/interviews/${id}/complete`, data)
}

export const addInterviewNote = (id, data) => {
  return request.put(`/interviews/${id}/note`, data)
}
