// Package server: AVS 的 HTTP请求 handlers
package server

import (
	"github.com/gin-gonic/gin"
	"goplus/shared/pkg/signature"
	"goplus/shared/pkg/types"
	"time"
)

type echoRequest struct {
	Data string `json:"data"`
}

type SignedOperatorResponse struct {
	Code        int                        `json:"code"`
	Message     string                     `json:"message"`
	Result      *types.SignedSecwareResult `json:"result"`
	SigOperator *types.HexBytes            `json:"sig_operator"`
}

func NewErrorOperatorResponse(code int, message string) SignedOperatorResponse {
	return SignedOperatorResponse{
		Code:        code,
		Message:     message,
		Result:      nil,
		SigOperator: nil,
	}
}

func (a *Server) handleTask(c *gin.Context) {
	taskStartTime := time.Now()
	task := types.SignedSecwareTask{}
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, NewErrorOperatorResponse(401, err.Error()))
		a.metricsIntf.IncTaskFailed()
		return
	}

	if !signature.VerifySignedSecwareTaskWithAddress(&task, a.config.AddressGateway) {
		c.JSON(400, NewErrorOperatorResponse(402, "bad SignedSecwareTask"))
		a.metricsIntf.IncTaskFailed()
		return
	}

	state, err := a.secwareManager.GetSecwareState(task.Task.SecwareId, task.Task.SecwareVersion)
	if err != nil {
		c.JSON(400, NewErrorOperatorResponse(403, err.Error()))
		a.metricsIntf.IncTaskFailed()
		return
	}

	task.Operator = a.config.AddressOperator[:]
	result, err := a.secwareAccessorIntf.HandleTask(state, &task)
	if err != nil {
		c.JSON(400, NewErrorOperatorResponse(404, err.Error()))
		a.metricsIntf.IncTaskFailed()
		return
	}
	signResult, err := signature.SignBLSOperatorResult(&result, a.config.BLSKeypair)
	if err != nil {
		c.JSON(400, NewErrorOperatorResponse(400, err.Error()))
		a.metricsIntf.IncTaskFailed()
		return
	}

	taskDuration := time.Since(taskStartTime).Seconds()
	a.metricsIntf.IncTaskHandled()
	a.metricsIntf.SetTaskDuration(taskDuration)

	response := SignedOperatorResponse{
		Code:        200,
		Message:     "ok",
		Result:      &signResult.Result,
		SigOperator: &signResult.SigOperator,
	}
	c.JSON(200, response)
}

func (a *Server) ping(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "pong"})
}
