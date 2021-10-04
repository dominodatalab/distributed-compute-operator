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
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;

/**
 * V1alpha1RayClusterSpecHeadEmptyDir
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2021-10-04T16:37:47.111Z[Etc/UTC]")
public class V1alpha1RayClusterSpecHeadEmptyDir {
  public static final String SERIALIZED_NAME_MEDIUM = "medium";
  @SerializedName(SERIALIZED_NAME_MEDIUM)
  private String medium;

  public static final String SERIALIZED_NAME_SIZE_LIMIT = "sizeLimit";
  @SerializedName(SERIALIZED_NAME_SIZE_LIMIT)
  private Object sizeLimit;


  public V1alpha1RayClusterSpecHeadEmptyDir medium(String medium) {
    
    this.medium = medium;
    return this;
  }

   /**
   * Get medium
   * @return medium
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public String getMedium() {
    return medium;
  }


  public void setMedium(String medium) {
    this.medium = medium;
  }


  public V1alpha1RayClusterSpecHeadEmptyDir sizeLimit(Object sizeLimit) {
    
    this.sizeLimit = sizeLimit;
    return this;
  }

   /**
   * Get sizeLimit
   * @return sizeLimit
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public Object getSizeLimit() {
    return sizeLimit;
  }


  public void setSizeLimit(Object sizeLimit) {
    this.sizeLimit = sizeLimit;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    V1alpha1RayClusterSpecHeadEmptyDir v1alpha1RayClusterSpecHeadEmptyDir = (V1alpha1RayClusterSpecHeadEmptyDir) o;
    return Objects.equals(this.medium, v1alpha1RayClusterSpecHeadEmptyDir.medium) &&
        Objects.equals(this.sizeLimit, v1alpha1RayClusterSpecHeadEmptyDir.sizeLimit);
  }

  @Override
  public int hashCode() {
    return Objects.hash(medium, sizeLimit);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class V1alpha1RayClusterSpecHeadEmptyDir {\n");
    sb.append("    medium: ").append(toIndentedString(medium)).append("\n");
    sb.append("    sizeLimit: ").append(toIndentedString(sizeLimit)).append("\n");
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

