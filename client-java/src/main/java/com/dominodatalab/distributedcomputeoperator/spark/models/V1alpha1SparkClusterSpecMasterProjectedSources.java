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


package com.dominodatalab.distributedcomputeoperator.spark.models;

import java.util.Objects;
import java.util.Arrays;
import com.dominodatalab.distributedcomputeoperator.spark.models.V1alpha1SparkClusterSpecMasterProjectedConfigMap;
import com.dominodatalab.distributedcomputeoperator.spark.models.V1alpha1SparkClusterSpecMasterProjectedDownwardAPI;
import com.dominodatalab.distributedcomputeoperator.spark.models.V1alpha1SparkClusterSpecMasterProjectedServiceAccountToken;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;

/**
 * V1alpha1SparkClusterSpecMasterProjectedSources
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2021-10-04T16:38:59.380Z[Etc/UTC]")
public class V1alpha1SparkClusterSpecMasterProjectedSources {
  public static final String SERIALIZED_NAME_CONFIG_MAP = "configMap";
  @SerializedName(SERIALIZED_NAME_CONFIG_MAP)
  private V1alpha1SparkClusterSpecMasterProjectedConfigMap configMap;

  public static final String SERIALIZED_NAME_DOWNWARD_A_P_I = "downwardAPI";
  @SerializedName(SERIALIZED_NAME_DOWNWARD_A_P_I)
  private V1alpha1SparkClusterSpecMasterProjectedDownwardAPI downwardAPI;

  public static final String SERIALIZED_NAME_SECRET = "secret";
  @SerializedName(SERIALIZED_NAME_SECRET)
  private V1alpha1SparkClusterSpecMasterProjectedConfigMap secret;

  public static final String SERIALIZED_NAME_SERVICE_ACCOUNT_TOKEN = "serviceAccountToken";
  @SerializedName(SERIALIZED_NAME_SERVICE_ACCOUNT_TOKEN)
  private V1alpha1SparkClusterSpecMasterProjectedServiceAccountToken serviceAccountToken;


  public V1alpha1SparkClusterSpecMasterProjectedSources configMap(V1alpha1SparkClusterSpecMasterProjectedConfigMap configMap) {
    
    this.configMap = configMap;
    return this;
  }

   /**
   * Get configMap
   * @return configMap
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1SparkClusterSpecMasterProjectedConfigMap getConfigMap() {
    return configMap;
  }


  public void setConfigMap(V1alpha1SparkClusterSpecMasterProjectedConfigMap configMap) {
    this.configMap = configMap;
  }


  public V1alpha1SparkClusterSpecMasterProjectedSources downwardAPI(V1alpha1SparkClusterSpecMasterProjectedDownwardAPI downwardAPI) {
    
    this.downwardAPI = downwardAPI;
    return this;
  }

   /**
   * Get downwardAPI
   * @return downwardAPI
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1SparkClusterSpecMasterProjectedDownwardAPI getDownwardAPI() {
    return downwardAPI;
  }


  public void setDownwardAPI(V1alpha1SparkClusterSpecMasterProjectedDownwardAPI downwardAPI) {
    this.downwardAPI = downwardAPI;
  }


  public V1alpha1SparkClusterSpecMasterProjectedSources secret(V1alpha1SparkClusterSpecMasterProjectedConfigMap secret) {
    
    this.secret = secret;
    return this;
  }

   /**
   * Get secret
   * @return secret
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1SparkClusterSpecMasterProjectedConfigMap getSecret() {
    return secret;
  }


  public void setSecret(V1alpha1SparkClusterSpecMasterProjectedConfigMap secret) {
    this.secret = secret;
  }


  public V1alpha1SparkClusterSpecMasterProjectedSources serviceAccountToken(V1alpha1SparkClusterSpecMasterProjectedServiceAccountToken serviceAccountToken) {
    
    this.serviceAccountToken = serviceAccountToken;
    return this;
  }

   /**
   * Get serviceAccountToken
   * @return serviceAccountToken
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1SparkClusterSpecMasterProjectedServiceAccountToken getServiceAccountToken() {
    return serviceAccountToken;
  }


  public void setServiceAccountToken(V1alpha1SparkClusterSpecMasterProjectedServiceAccountToken serviceAccountToken) {
    this.serviceAccountToken = serviceAccountToken;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    V1alpha1SparkClusterSpecMasterProjectedSources v1alpha1SparkClusterSpecMasterProjectedSources = (V1alpha1SparkClusterSpecMasterProjectedSources) o;
    return Objects.equals(this.configMap, v1alpha1SparkClusterSpecMasterProjectedSources.configMap) &&
        Objects.equals(this.downwardAPI, v1alpha1SparkClusterSpecMasterProjectedSources.downwardAPI) &&
        Objects.equals(this.secret, v1alpha1SparkClusterSpecMasterProjectedSources.secret) &&
        Objects.equals(this.serviceAccountToken, v1alpha1SparkClusterSpecMasterProjectedSources.serviceAccountToken);
  }

  @Override
  public int hashCode() {
    return Objects.hash(configMap, downwardAPI, secret, serviceAccountToken);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class V1alpha1SparkClusterSpecMasterProjectedSources {\n");
    sb.append("    configMap: ").append(toIndentedString(configMap)).append("\n");
    sb.append("    downwardAPI: ").append(toIndentedString(downwardAPI)).append("\n");
    sb.append("    secret: ").append(toIndentedString(secret)).append("\n");
    sb.append("    serviceAccountToken: ").append(toIndentedString(serviceAccountToken)).append("\n");
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
