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


package com.dominodatalab.distributedcomputeoperator.dask.models;

import java.util.Objects;
import java.util.Arrays;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;

/**
 * V1alpha1DaskClusterSpecValueFromFieldRef
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2021-10-04T16:36:51.907Z[Etc/UTC]")
public class V1alpha1DaskClusterSpecValueFromFieldRef {
  public static final String SERIALIZED_NAME_API_VERSION = "apiVersion";
  @SerializedName(SERIALIZED_NAME_API_VERSION)
  private String apiVersion;

  public static final String SERIALIZED_NAME_FIELD_PATH = "fieldPath";
  @SerializedName(SERIALIZED_NAME_FIELD_PATH)
  private String fieldPath;


  public V1alpha1DaskClusterSpecValueFromFieldRef apiVersion(String apiVersion) {
    
    this.apiVersion = apiVersion;
    return this;
  }

   /**
   * Get apiVersion
   * @return apiVersion
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public String getApiVersion() {
    return apiVersion;
  }


  public void setApiVersion(String apiVersion) {
    this.apiVersion = apiVersion;
  }


  public V1alpha1DaskClusterSpecValueFromFieldRef fieldPath(String fieldPath) {
    
    this.fieldPath = fieldPath;
    return this;
  }

   /**
   * Get fieldPath
   * @return fieldPath
  **/
  @ApiModelProperty(required = true, value = "")

  public String getFieldPath() {
    return fieldPath;
  }


  public void setFieldPath(String fieldPath) {
    this.fieldPath = fieldPath;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    V1alpha1DaskClusterSpecValueFromFieldRef v1alpha1DaskClusterSpecValueFromFieldRef = (V1alpha1DaskClusterSpecValueFromFieldRef) o;
    return Objects.equals(this.apiVersion, v1alpha1DaskClusterSpecValueFromFieldRef.apiVersion) &&
        Objects.equals(this.fieldPath, v1alpha1DaskClusterSpecValueFromFieldRef.fieldPath);
  }

  @Override
  public int hashCode() {
    return Objects.hash(apiVersion, fieldPath);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class V1alpha1DaskClusterSpecValueFromFieldRef {\n");
    sb.append("    apiVersion: ").append(toIndentedString(apiVersion)).append("\n");
    sb.append("    fieldPath: ").append(toIndentedString(fieldPath)).append("\n");
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

