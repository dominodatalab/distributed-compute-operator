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
 * V1alpha1DaskClusterSpecSchedulerSpecDataSource
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2021-10-04T16:36:51.907Z[Etc/UTC]")
public class V1alpha1DaskClusterSpecSchedulerSpecDataSource {
  public static final String SERIALIZED_NAME_API_GROUP = "apiGroup";
  @SerializedName(SERIALIZED_NAME_API_GROUP)
  private String apiGroup;

  public static final String SERIALIZED_NAME_KIND = "kind";
  @SerializedName(SERIALIZED_NAME_KIND)
  private String kind;

  public static final String SERIALIZED_NAME_NAME = "name";
  @SerializedName(SERIALIZED_NAME_NAME)
  private String name;


  public V1alpha1DaskClusterSpecSchedulerSpecDataSource apiGroup(String apiGroup) {
    
    this.apiGroup = apiGroup;
    return this;
  }

   /**
   * Get apiGroup
   * @return apiGroup
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public String getApiGroup() {
    return apiGroup;
  }


  public void setApiGroup(String apiGroup) {
    this.apiGroup = apiGroup;
  }


  public V1alpha1DaskClusterSpecSchedulerSpecDataSource kind(String kind) {
    
    this.kind = kind;
    return this;
  }

   /**
   * Get kind
   * @return kind
  **/
  @ApiModelProperty(required = true, value = "")

  public String getKind() {
    return kind;
  }


  public void setKind(String kind) {
    this.kind = kind;
  }


  public V1alpha1DaskClusterSpecSchedulerSpecDataSource name(String name) {
    
    this.name = name;
    return this;
  }

   /**
   * Get name
   * @return name
  **/
  @ApiModelProperty(required = true, value = "")

  public String getName() {
    return name;
  }


  public void setName(String name) {
    this.name = name;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    V1alpha1DaskClusterSpecSchedulerSpecDataSource v1alpha1DaskClusterSpecSchedulerSpecDataSource = (V1alpha1DaskClusterSpecSchedulerSpecDataSource) o;
    return Objects.equals(this.apiGroup, v1alpha1DaskClusterSpecSchedulerSpecDataSource.apiGroup) &&
        Objects.equals(this.kind, v1alpha1DaskClusterSpecSchedulerSpecDataSource.kind) &&
        Objects.equals(this.name, v1alpha1DaskClusterSpecSchedulerSpecDataSource.name);
  }

  @Override
  public int hashCode() {
    return Objects.hash(apiGroup, kind, name);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class V1alpha1DaskClusterSpecSchedulerSpecDataSource {\n");
    sb.append("    apiGroup: ").append(toIndentedString(apiGroup)).append("\n");
    sb.append("    kind: ").append(toIndentedString(kind)).append("\n");
    sb.append("    name: ").append(toIndentedString(name)).append("\n");
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
