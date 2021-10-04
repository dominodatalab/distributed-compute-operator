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
import com.dominodatalab.distributedcomputeoperator.ray.models.V1alpha1RayClusterSpecHeadProjectedConfigMap;
import com.dominodatalab.distributedcomputeoperator.ray.models.V1alpha1RayClusterSpecHeadProjectedDownwardAPI;
import com.dominodatalab.distributedcomputeoperator.ray.models.V1alpha1RayClusterSpecHeadProjectedServiceAccountToken;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;

/**
 * V1alpha1RayClusterSpecHeadProjectedSources
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2021-10-04T16:37:47.111Z[Etc/UTC]")
public class V1alpha1RayClusterSpecHeadProjectedSources {
  public static final String SERIALIZED_NAME_CONFIG_MAP = "configMap";
  @SerializedName(SERIALIZED_NAME_CONFIG_MAP)
  private V1alpha1RayClusterSpecHeadProjectedConfigMap configMap;

  public static final String SERIALIZED_NAME_DOWNWARD_A_P_I = "downwardAPI";
  @SerializedName(SERIALIZED_NAME_DOWNWARD_A_P_I)
  private V1alpha1RayClusterSpecHeadProjectedDownwardAPI downwardAPI;

  public static final String SERIALIZED_NAME_SECRET = "secret";
  @SerializedName(SERIALIZED_NAME_SECRET)
  private V1alpha1RayClusterSpecHeadProjectedConfigMap secret;

  public static final String SERIALIZED_NAME_SERVICE_ACCOUNT_TOKEN = "serviceAccountToken";
  @SerializedName(SERIALIZED_NAME_SERVICE_ACCOUNT_TOKEN)
  private V1alpha1RayClusterSpecHeadProjectedServiceAccountToken serviceAccountToken;


  public V1alpha1RayClusterSpecHeadProjectedSources configMap(V1alpha1RayClusterSpecHeadProjectedConfigMap configMap) {
    
    this.configMap = configMap;
    return this;
  }

   /**
   * Get configMap
   * @return configMap
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1RayClusterSpecHeadProjectedConfigMap getConfigMap() {
    return configMap;
  }


  public void setConfigMap(V1alpha1RayClusterSpecHeadProjectedConfigMap configMap) {
    this.configMap = configMap;
  }


  public V1alpha1RayClusterSpecHeadProjectedSources downwardAPI(V1alpha1RayClusterSpecHeadProjectedDownwardAPI downwardAPI) {
    
    this.downwardAPI = downwardAPI;
    return this;
  }

   /**
   * Get downwardAPI
   * @return downwardAPI
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1RayClusterSpecHeadProjectedDownwardAPI getDownwardAPI() {
    return downwardAPI;
  }


  public void setDownwardAPI(V1alpha1RayClusterSpecHeadProjectedDownwardAPI downwardAPI) {
    this.downwardAPI = downwardAPI;
  }


  public V1alpha1RayClusterSpecHeadProjectedSources secret(V1alpha1RayClusterSpecHeadProjectedConfigMap secret) {
    
    this.secret = secret;
    return this;
  }

   /**
   * Get secret
   * @return secret
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1RayClusterSpecHeadProjectedConfigMap getSecret() {
    return secret;
  }


  public void setSecret(V1alpha1RayClusterSpecHeadProjectedConfigMap secret) {
    this.secret = secret;
  }


  public V1alpha1RayClusterSpecHeadProjectedSources serviceAccountToken(V1alpha1RayClusterSpecHeadProjectedServiceAccountToken serviceAccountToken) {
    
    this.serviceAccountToken = serviceAccountToken;
    return this;
  }

   /**
   * Get serviceAccountToken
   * @return serviceAccountToken
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1RayClusterSpecHeadProjectedServiceAccountToken getServiceAccountToken() {
    return serviceAccountToken;
  }


  public void setServiceAccountToken(V1alpha1RayClusterSpecHeadProjectedServiceAccountToken serviceAccountToken) {
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
    V1alpha1RayClusterSpecHeadProjectedSources v1alpha1RayClusterSpecHeadProjectedSources = (V1alpha1RayClusterSpecHeadProjectedSources) o;
    return Objects.equals(this.configMap, v1alpha1RayClusterSpecHeadProjectedSources.configMap) &&
        Objects.equals(this.downwardAPI, v1alpha1RayClusterSpecHeadProjectedSources.downwardAPI) &&
        Objects.equals(this.secret, v1alpha1RayClusterSpecHeadProjectedSources.secret) &&
        Objects.equals(this.serviceAccountToken, v1alpha1RayClusterSpecHeadProjectedSources.serviceAccountToken);
  }

  @Override
  public int hashCode() {
    return Objects.hash(configMap, downwardAPI, secret, serviceAccountToken);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class V1alpha1RayClusterSpecHeadProjectedSources {\n");
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

