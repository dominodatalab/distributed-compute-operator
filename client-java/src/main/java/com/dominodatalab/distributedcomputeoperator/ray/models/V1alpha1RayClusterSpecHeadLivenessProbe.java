/*
 * Kubernetes
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: v1.18.2
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package com.dominodatalab.distributedcomputeoperator.ray.models;

import java.util.Objects;
import java.util.Arrays;
import com.dominodatalab.distributedcomputeoperator.ray.models.V1alpha1RayClusterSpecHeadLifecyclePostStartExec;
import com.dominodatalab.distributedcomputeoperator.ray.models.V1alpha1RayClusterSpecHeadLifecyclePostStartHttpGet;
import com.dominodatalab.distributedcomputeoperator.ray.models.V1alpha1RayClusterSpecHeadLifecyclePostStartTcpSocket;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;

/**
 * V1alpha1RayClusterSpecHeadLivenessProbe
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2021-10-04T16:37:47.111Z[Etc/UTC]")
public class V1alpha1RayClusterSpecHeadLivenessProbe {
  public static final String SERIALIZED_NAME_EXEC = "exec";
  @SerializedName(SERIALIZED_NAME_EXEC)
  private V1alpha1RayClusterSpecHeadLifecyclePostStartExec exec;

  public static final String SERIALIZED_NAME_FAILURE_THRESHOLD = "failureThreshold";
  @SerializedName(SERIALIZED_NAME_FAILURE_THRESHOLD)
  private Integer failureThreshold;

  public static final String SERIALIZED_NAME_HTTP_GET = "httpGet";
  @SerializedName(SERIALIZED_NAME_HTTP_GET)
  private V1alpha1RayClusterSpecHeadLifecyclePostStartHttpGet httpGet;

  public static final String SERIALIZED_NAME_INITIAL_DELAY_SECONDS = "initialDelaySeconds";
  @SerializedName(SERIALIZED_NAME_INITIAL_DELAY_SECONDS)
  private Integer initialDelaySeconds;

  public static final String SERIALIZED_NAME_PERIOD_SECONDS = "periodSeconds";
  @SerializedName(SERIALIZED_NAME_PERIOD_SECONDS)
  private Integer periodSeconds;

  public static final String SERIALIZED_NAME_SUCCESS_THRESHOLD = "successThreshold";
  @SerializedName(SERIALIZED_NAME_SUCCESS_THRESHOLD)
  private Integer successThreshold;

  public static final String SERIALIZED_NAME_TCP_SOCKET = "tcpSocket";
  @SerializedName(SERIALIZED_NAME_TCP_SOCKET)
  private V1alpha1RayClusterSpecHeadLifecyclePostStartTcpSocket tcpSocket;

  public static final String SERIALIZED_NAME_TERMINATION_GRACE_PERIOD_SECONDS = "terminationGracePeriodSeconds";
  @SerializedName(SERIALIZED_NAME_TERMINATION_GRACE_PERIOD_SECONDS)
  private Long terminationGracePeriodSeconds;

  public static final String SERIALIZED_NAME_TIMEOUT_SECONDS = "timeoutSeconds";
  @SerializedName(SERIALIZED_NAME_TIMEOUT_SECONDS)
  private Integer timeoutSeconds;


  public V1alpha1RayClusterSpecHeadLivenessProbe exec(V1alpha1RayClusterSpecHeadLifecyclePostStartExec exec) {
    
    this.exec = exec;
    return this;
  }

   /**
   * Get exec
   * @return exec
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1RayClusterSpecHeadLifecyclePostStartExec getExec() {
    return exec;
  }


  public void setExec(V1alpha1RayClusterSpecHeadLifecyclePostStartExec exec) {
    this.exec = exec;
  }


  public V1alpha1RayClusterSpecHeadLivenessProbe failureThreshold(Integer failureThreshold) {
    
    this.failureThreshold = failureThreshold;
    return this;
  }

   /**
   * Get failureThreshold
   * @return failureThreshold
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public Integer getFailureThreshold() {
    return failureThreshold;
  }


  public void setFailureThreshold(Integer failureThreshold) {
    this.failureThreshold = failureThreshold;
  }


  public V1alpha1RayClusterSpecHeadLivenessProbe httpGet(V1alpha1RayClusterSpecHeadLifecyclePostStartHttpGet httpGet) {
    
    this.httpGet = httpGet;
    return this;
  }

   /**
   * Get httpGet
   * @return httpGet
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1RayClusterSpecHeadLifecyclePostStartHttpGet getHttpGet() {
    return httpGet;
  }


  public void setHttpGet(V1alpha1RayClusterSpecHeadLifecyclePostStartHttpGet httpGet) {
    this.httpGet = httpGet;
  }


  public V1alpha1RayClusterSpecHeadLivenessProbe initialDelaySeconds(Integer initialDelaySeconds) {
    
    this.initialDelaySeconds = initialDelaySeconds;
    return this;
  }

   /**
   * Get initialDelaySeconds
   * @return initialDelaySeconds
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public Integer getInitialDelaySeconds() {
    return initialDelaySeconds;
  }


  public void setInitialDelaySeconds(Integer initialDelaySeconds) {
    this.initialDelaySeconds = initialDelaySeconds;
  }


  public V1alpha1RayClusterSpecHeadLivenessProbe periodSeconds(Integer periodSeconds) {
    
    this.periodSeconds = periodSeconds;
    return this;
  }

   /**
   * Get periodSeconds
   * @return periodSeconds
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public Integer getPeriodSeconds() {
    return periodSeconds;
  }


  public void setPeriodSeconds(Integer periodSeconds) {
    this.periodSeconds = periodSeconds;
  }


  public V1alpha1RayClusterSpecHeadLivenessProbe successThreshold(Integer successThreshold) {
    
    this.successThreshold = successThreshold;
    return this;
  }

   /**
   * Get successThreshold
   * @return successThreshold
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public Integer getSuccessThreshold() {
    return successThreshold;
  }


  public void setSuccessThreshold(Integer successThreshold) {
    this.successThreshold = successThreshold;
  }


  public V1alpha1RayClusterSpecHeadLivenessProbe tcpSocket(V1alpha1RayClusterSpecHeadLifecyclePostStartTcpSocket tcpSocket) {
    
    this.tcpSocket = tcpSocket;
    return this;
  }

   /**
   * Get tcpSocket
   * @return tcpSocket
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1RayClusterSpecHeadLifecyclePostStartTcpSocket getTcpSocket() {
    return tcpSocket;
  }


  public void setTcpSocket(V1alpha1RayClusterSpecHeadLifecyclePostStartTcpSocket tcpSocket) {
    this.tcpSocket = tcpSocket;
  }


  public V1alpha1RayClusterSpecHeadLivenessProbe terminationGracePeriodSeconds(Long terminationGracePeriodSeconds) {
    
    this.terminationGracePeriodSeconds = terminationGracePeriodSeconds;
    return this;
  }

   /**
   * Get terminationGracePeriodSeconds
   * @return terminationGracePeriodSeconds
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public Long getTerminationGracePeriodSeconds() {
    return terminationGracePeriodSeconds;
  }


  public void setTerminationGracePeriodSeconds(Long terminationGracePeriodSeconds) {
    this.terminationGracePeriodSeconds = terminationGracePeriodSeconds;
  }


  public V1alpha1RayClusterSpecHeadLivenessProbe timeoutSeconds(Integer timeoutSeconds) {
    
    this.timeoutSeconds = timeoutSeconds;
    return this;
  }

   /**
   * Get timeoutSeconds
   * @return timeoutSeconds
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public Integer getTimeoutSeconds() {
    return timeoutSeconds;
  }


  public void setTimeoutSeconds(Integer timeoutSeconds) {
    this.timeoutSeconds = timeoutSeconds;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    V1alpha1RayClusterSpecHeadLivenessProbe v1alpha1RayClusterSpecHeadLivenessProbe = (V1alpha1RayClusterSpecHeadLivenessProbe) o;
    return Objects.equals(this.exec, v1alpha1RayClusterSpecHeadLivenessProbe.exec) &&
        Objects.equals(this.failureThreshold, v1alpha1RayClusterSpecHeadLivenessProbe.failureThreshold) &&
        Objects.equals(this.httpGet, v1alpha1RayClusterSpecHeadLivenessProbe.httpGet) &&
        Objects.equals(this.initialDelaySeconds, v1alpha1RayClusterSpecHeadLivenessProbe.initialDelaySeconds) &&
        Objects.equals(this.periodSeconds, v1alpha1RayClusterSpecHeadLivenessProbe.periodSeconds) &&
        Objects.equals(this.successThreshold, v1alpha1RayClusterSpecHeadLivenessProbe.successThreshold) &&
        Objects.equals(this.tcpSocket, v1alpha1RayClusterSpecHeadLivenessProbe.tcpSocket) &&
        Objects.equals(this.terminationGracePeriodSeconds, v1alpha1RayClusterSpecHeadLivenessProbe.terminationGracePeriodSeconds) &&
        Objects.equals(this.timeoutSeconds, v1alpha1RayClusterSpecHeadLivenessProbe.timeoutSeconds);
  }

  @Override
  public int hashCode() {
    return Objects.hash(exec, failureThreshold, httpGet, initialDelaySeconds, periodSeconds, successThreshold, tcpSocket, terminationGracePeriodSeconds, timeoutSeconds);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class V1alpha1RayClusterSpecHeadLivenessProbe {\n");
    sb.append("    exec: ").append(toIndentedString(exec)).append("\n");
    sb.append("    failureThreshold: ").append(toIndentedString(failureThreshold)).append("\n");
    sb.append("    httpGet: ").append(toIndentedString(httpGet)).append("\n");
    sb.append("    initialDelaySeconds: ").append(toIndentedString(initialDelaySeconds)).append("\n");
    sb.append("    periodSeconds: ").append(toIndentedString(periodSeconds)).append("\n");
    sb.append("    successThreshold: ").append(toIndentedString(successThreshold)).append("\n");
    sb.append("    tcpSocket: ").append(toIndentedString(tcpSocket)).append("\n");
    sb.append("    terminationGracePeriodSeconds: ").append(toIndentedString(terminationGracePeriodSeconds)).append("\n");
    sb.append("    timeoutSeconds: ").append(toIndentedString(timeoutSeconds)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }

}

